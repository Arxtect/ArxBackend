package controllers

import (
	"github.com/Arxtect/ArxBackend/golangp/apps/arx_center/models"
	"github.com/Arxtect/ArxBackend/golangp/common/constants"
	"github.com/Arxtect/ArxBackend/golangp/common/logger"
	openai_config "github.com/Arxtect/ArxBackend/golangp/common/openai-config"
	"github.com/Arxtect/ArxBackend/golangp/config"
	"context"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkoukk/tiktoken-go"
	"github.com/toheart/functrace"

	"github.com/gin-gonic/gin"
	gogpt "github.com/sashabaranov/go-openai"
	"golang.org/x/net/proxy"
)

type BaseController struct {
}

func (*BaseController) ResponseJson(ctx *gin.Context, code int, errorMsg string, data interface{}) {
	defer functrace.Trace([]interface {
	}{ctx, code, errorMsg, data})()

	ctx.JSON(code, gin.H{
		"code":     code,
		"errorMsg": errorMsg,
		"data":     data,
	})
	ctx.Abort()
}

func (*BaseController) ResponseData(ctx *gin.Context, code int, contentType string, data []byte) {
	defer functrace.Trace([]interface {
	}{ctx, code, contentType, data})()
	ctx.Data(code, contentType, data)
	ctx.Abort()
}

type ChatController struct {
	BaseController
	cs *CreditSystem
}

func NewChatController(creditSystem *CreditSystem) ChatController {
	defer functrace.Trace([]interface {
	}{creditSystem})()
	return ChatController{cs: creditSystem}
}

func NumTokens(text string) int {
	defer functrace.Trace([]interface {
	}{text})()
	encoding := "cl100k_base"

	tke, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		err = fmt.Errorf("GetEncoding: %v", err)
		return -1
	}

	token := tke.Encode(text, nil, nil)

	numTokens := len(token)
	return numTokens
}

func (c *ChatController) CompletionWithModelInfo(ctx *gin.Context) {
	defer functrace.Trace([]interface {
	}{c, ctx})()
	var request gogpt.ChatCompletionRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	logger.Info("Starting a request with model %s", request.Model)
	if len(request.Messages) == 0 {
		c.ResponseJson(ctx, http.StatusBadRequest, "request messages required", nil)
		return
	}

	configCopy := config.Env

	cnf := openai_config.OpenAIConfiguration{
		ApiKey:        configCopy.ApiKey,
		ApiURL:        configCopy.ApiURL,
		Listen:        configCopy.Listen,
		Proxy:         configCopy.Proxy,
		AdminEmail:    configCopy.AdminEmail,
		AdminPassword: configCopy.AdminPassword,
	}

	gptConfig := gogpt.DefaultConfig(cnf.ApiKey)

	if cnf.Proxy != "" {
		transport := &http.Transport{}

		if strings.HasPrefix(cnf.Proxy, "socks5h://") {

			dialContext, err := newDialContext(cnf.Proxy[10:])
			if err != nil {
				panic(err)
			}
			transport.DialContext = dialContext
		} else {

			proxyUrl, err := url.Parse(cnf.Proxy)
			if err != nil {
				panic(err)
			}
			transport.Proxy = http.ProxyURL(proxyUrl)
		}

		gptConfig.HTTPClient = &http.Client{
			Transport: transport,
		}

	}

	if cnf.ApiURL != "" {
		gptConfig.BaseURL = cnf.ApiURL
	}

	client := gogpt.NewClientWithConfig(gptConfig)
	if request.Model == "" {
		logger.Danger("request model is empty")
		c.ResponseJson(ctx, http.StatusBadRequest, "request model is empty", nil)
	}

	currentUser := ctx.MustGet("currentUser").(models.User)
	if !c.preCheckBalance(ctx, currentUser.Balance) {
		return
	}

	if request.Model == gogpt.GPT3Dot5Turbo0301 || request.Model == gogpt.GPT4 || request.Model == gogpt.
		GPT40314 || request.
		Model == gogpt.
		GPT3Dot5Turbo {
		if request.Stream {
			logger.Info("stream request started")
			stream, err := client.CreateChatCompletionStream(ctx, request)
			tokenPromptCount := 0
			for _, msg := range request.Messages {
				tokenPromptCount += NumTokens(msg.Content)
			}
			if err != nil {
				c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
				return
			}

			chanStream := make(chan string, 10)
			res := ""

			go func() {
				for {
					nextResp, err := stream.Recv()
					if err == io.EOF {
						stream.Close()
						tokenCompletionCount := NumTokens(res)
						totalTokenCount := tokenPromptCount + tokenCompletionCount
						usage := gogpt.Usage{PromptTokens: tokenPromptCount, CompletionTokens: tokenCompletionCount, TotalTokens: totalTokenCount}
						cost := GetModelChargeInChineseCents(request.Model, usage)
						chanStream <- fmt.Sprintf("[TotalCredit: %d]", cost)
						c.chargeUserFromBalance(currentUser.Email, cost)
						return
					} else if err != nil {
						c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
						return
					} else {
						chanStream <- nextResp.Choices[0].Delta.Content
						res += nextResp.Choices[0].Delta.Content
					}
				}

			}()

			ctx.Stream(func(w io.Writer) bool {
				if msg, ok := <-chanStream; ok {
					if !strings.HasPrefix(msg, "[TotalCredit:") {
						_, err := w.Write([]byte(msg))
						if err != nil {
							logger.Warning(err.Error())
							return false
						}
						return true
					} else {
						_, err := w.Write([]byte(msg))
						if err != nil {
							logger.Warning(err.Error())
							return false
						}
						close(chanStream)
						return true
					}
				}
				return false
			})

		} else {
			resp, err := client.CreateChatCompletion(ctx, request)
			if err != nil {
				c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
				return
			}
			cost := GetModelChargeInChineseCents(request.Model, resp.Usage)
			c.chargeUserFromBalance(currentUser.Email, cost)

			c.ResponseJson(ctx, http.StatusOK, "", gin.H{
				"reply":       resp.Choices[0].Message.Content,
				"messages":    append(request.Messages, resp.Choices[0].Message),
				"totalCredit": cost,
			})
		}
	} else {
		prompt := ""
		for _, item := range request.Messages {
			prompt += item.Content + "/n"
		}
		prompt = strings.Trim(prompt, "/n")

		logger.Info("request prompt is %s", prompt)
		req := gogpt.CompletionRequest{
			Model:            request.Model,
			MaxTokens:        request.MaxTokens,
			TopP:             request.TopP,
			FrequencyPenalty: request.FrequencyPenalty,
			PresencePenalty:  request.PresencePenalty,
			Prompt:           prompt,
		}

		resp, err := client.CreateCompletion(ctx, req)

		cost := GetModelChargeInChineseCents(request.Model, resp.Usage)
		c.chargeUserFromBalance(currentUser.Email, cost)

		if err != nil {
			c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		c.ResponseJson(ctx, http.StatusOK, "", gin.H{
			"reply": resp.Choices[0].Text,
			"messages": append(request.Messages, gogpt.ChatCompletionMessage{
				Role:    "assistant",
				Content: resp.Choices[0].Text,
			}),
		})
	}

}

func (c *ChatController) preCheckBalance(ctx *gin.Context, balance int64) bool {
	defer functrace.Trace([]interface {
	}{c, ctx, balance})()
	if balance < 5 {
		c.ResponseJson(ctx, http.StatusBadRequest, "Insufficient balance", nil)
		return false
	}
	return true
}

func (c *ChatController) chargeUserFromBalance(email string, cost int64) {
	defer functrace.Trace([]interface {
	}{c, email, cost})()
	_, err := c.cs.UpdateBalanceByUserEmail(email, -1*cost)
	if err != nil {
		logger.Warning("charge user %s at cost %d failed", email, cost)
	} else {
		logger.Info("charge user %s at cost %d success", email, cost)
	}
}

func GetModelChargeInChineseCents(model string, usage gogpt.Usage) int64 {
	defer functrace.Trace([]interface {
	}{model, usage})()
	var chargeInChineseCents float64 = 0
	switch {
	case model == gogpt.GPT3Dot5Turbo0301 || model == gogpt.GPT3Dot5Turbo:
		chargeInChineseCents = (float64(usage.PromptTokens)*(constants.GPT3PromptCharge) + float64(usage.
			CompletionTokens)*(constants.GPT3CompletionCharge)) * constants.DollarToChineseCentsRate
		break
	case model == gogpt.GPT4 || model == gogpt.GPT40314:
		chargeInChineseCents = (float64(usage.PromptTokens)*(constants.GPT4PromptCharge) + float64(usage.
			CompletionTokens)*(constants.GPT4CompletionCharge)) * constants.DollarToChineseCentsRate
		break
	}
	return int64(math.Ceil(chargeInChineseCents))
}

type dialContextFunc func(ctx context.Context, network, address string) (net.Conn, error)

func newDialContext(socks5 string) (dialContextFunc, error) {
	defer functrace.Trace([]interface {
	}{socks5})()
	baseDialer := &net.Dialer{
		Timeout:   60 * time.Second,
		KeepAlive: 60 * time.Second,
	}

	if socks5 != "" {

		var auth *proxy.Auth = nil

		if strings.Contains(socks5, "@") {
			proxyInfo := strings.SplitN(socks5, "@", 2)
			proxyUser := strings.Split(proxyInfo[0], ":")
			if len(proxyUser) == 2 {
				auth = &proxy.Auth{
					User:     proxyUser[0],
					Password: proxyUser[1],
				}
			}
			socks5 = proxyInfo[1]
		}

		dialSocksProxy, err := proxy.SOCKS5("tcp", socks5, auth, baseDialer)
		if err != nil {
			return nil, err
		}

		contextDialer, ok := dialSocksProxy.(proxy.ContextDialer)
		if !ok {
			return nil, err
		}

		return contextDialer.DialContext, nil
	} else {
		return baseDialer.DialContext, nil
	}
}
