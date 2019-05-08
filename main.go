package main

import (
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var websocketUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	//启动WebSocket
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		conn, upgradeError := websocketUpgrader.Upgrade(writer, request, nil)
		if upgradeError != nil {
			log.Fatal("协议升级失败: ", upgradeError)
			return
		}
		defer conn.Close()

		messageType, messageBytes, readError := conn.ReadMessage()
		if readError != nil {
			log.Fatal("读取消息失败: ", readError)
			return
		}
		log.Printf("消息类型:%v , 消息内容:%s", messageType, string(messageBytes))

		//返回测试消息
		var buffer bytes.Buffer
		for i := 1; i <= 10000000; i++ {
			buffer.WriteString("字")
		}
		log.Println("字符长度:", len(buffer.String()))
		writeError := conn.WriteMessage(websocket.TextMessage, buffer.Bytes())
		if writeError != nil {
			log.Fatal("发送消息错误:", writeError)
		}

		//发送关闭消息
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "处理完毕"))
	})

	log.Println("WebSocket:", 8080)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("WebSocket失败:", err)
	}
}
