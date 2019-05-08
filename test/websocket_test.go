package test

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"testing"
	"time"
)

func TestWebSocket(t *testing.T) {
	log.Println("测试WebSocket")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8080", Path: "/"}
	log.Printf("连接到: %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			messageType, messageData, readError := c.ReadMessage()
			if readError != nil {
				if readError.Error() == "websocket: close 1000 (normal): 处理完毕" {
					log.Println("正常关闭")
				} else {
					log.Println("读取消息错误:", readError)
				}
				return
			}
			log.Printf("收到消息-类型:%v , 字符长度:%d", messageType, len(string(messageData)))
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			log.Println("发送消息:", t.String())
			writeError := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("发送消息失败:", writeError)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
