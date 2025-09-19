package main

import (
	// "github.com/Ating56/websocket-plugin"

	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	wp "websocket-plugin"
)

var Redis = wp.Redis{
	Host:     "127.0.0.1:6379",
	Password: "",
	DB:       0,
}
var MongoDB = wp.MongoDB{
	DataSource: "mongodb://root:student@127.0.0.1:27017",
	DataBase:   "test",
	Collection: "websocket",
}

// 消息结构体自定义，可包含前端传来的字段，额外再加发送者的id
type Message struct {
	MsgId     string
	SenderId  string
	ReceiveId string
	Content   string
	Time      string
}

func main() {
	wp.NewConfig(Redis, MongoDB)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		clientId := r.Header.Get("Sec-WebSocket-Protocol")
		wp.SetConnect(w, r, clientId)
	})
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		var target Message
		clientId := r.Header.Get("ClientId")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(body, &target)
		if err != nil {
			return
		}

		target.MsgId = fmt.Sprintf("msg%d", time.Now().Unix()) // 消息唯一id，可自定义
		target.SenderId = clientId
		target.Time = time.Now().Format("2006-01-02 15:04:05")

		dataSendToWs, err := json.Marshal(target)
		if err != nil {
			log.Println("json转化失败")
			return
		}
		fmt.Println("body", target.ReceiveId, target.Content, clientId)

		wp.SendToServer(clientId, target.ReceiveId, dataSendToWs)
	})
	http.Handle("/", http.FileServer(http.Dir("client")))
	http.ListenAndServe(":8080", nil)
}
