package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/toheart/functrace"
	"golang.org/x/net/websocket"
)

type subscribeRequest struct {
	RoomID   string `json:"roomId"`
	UserID   string `json:"userId"`
	RoleID   string `json:"roleID"`
	UserName string `json:"userName"`
}

func HandlerWs(c *gin.Context) {
	defer functrace.Trace([]interface {
	}{c})()
	ws := c.Writer
	req := c.Request

	websocket.Handler(func(conn *websocket.Conn) {
		defer conn.Close()

		for {
			var message string

			if err := websocket.Message.Receive(conn, &message); err != nil {

				fmt.Printf("Error reading from WebSocket: %v\n", err)
				break
			}

			var subReq subscribeRequest
			if err := json.Unmarshal([]byte(message), &subReq); err != nil {
				fmt.Printf("Error unmarshalling subscribe request: %v\n", err)
				continue
			}

			subscribers, err := ssList.GetSubscribers(subReq.RoomID)
			if err != nil {
				fmt.Printf("Error getting subscribers for room: %v\n", err)
				continue
			}

			for _, subscriber := range subscribers {
				if subscriber.Connection != conn {

					if err := websocket.Message.Send(subscriber.Connection, message); err != nil {
						fmt.Printf("Error sending message to subscriber: %v\n", err)

						continue
					}
				}
			}
		}
	}).ServeHTTP(ws, req)
}
