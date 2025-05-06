package openai_config

import (
	"encoding/json"
	"github.com/toheart/functrace"
	"log"
	"os"
	"sync"
)

type OpenAIConfiguration struct {
	ApiKey string `json:"api_key"`

	ApiURL string `json:"api_url"`

	Listen string `json:"listen"`

	Proxy         string   `json:"proxy"`
	AdminEmail    []string `json:"admin_email"`
	AdminPassword string   `json:"admin_password"`
}

var config *OpenAIConfiguration
var once sync.Once

func LoadOpenAIConfig() *OpenAIConfiguration {
	defer functrace.Trace([]interface {
	}{})()
	once.Do(func() {

		config = &OpenAIConfiguration{
			ApiURL: "",
			Listen: "",
		}

		_, err := os.Stat(CLI.Config)
		if err == nil {
			f, err := os.Open(CLI.Config)
			if err != nil {
				log.Fatalf("open openai-config err: %v", err)
				return
			}
			defer f.Close()
			encoder := json.NewDecoder(f)
			err = encoder.Decode(config)
			if err != nil {
				log.Fatalf("decode openai-config err: %v", err)
				return
			}
		}
	})

	return config
}
