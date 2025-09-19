package websocketplugin

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ClientId   string // 等效于发送用户Id
	Conn       *websocket.Conn
	RemoteAddr string
}

var GlobalRecv sync.Map

func SendToServer(targetId, text string) {
	msg := map[string][]byte{
		"targetId": []byte(targetId),
		"message":  []byte(text),
	}
	fmt.Println("SendToServer targetId:", targetId, "text:", msg)
	if recvChan, ok := GlobalRecv.Load(targetId); ok {
		recvChan.(chan map[string][]byte) <- msg
	}
}

// sendToClient 持续发送消息到客户端
// 工作流程：
// 1. 创建30秒间隔的心跳ticker
// 2. 监听消息通道和心跳事件
// 3. 处理消息发送或Ping心跳包
// 注意：
// - 设置10秒写超时防止阻塞
// - 通道关闭时退出循环
func (c *Client) SendToClient() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		targetChan, ok := GlobalRecv.Load(c.ClientId)
		fmt.Println("targetChan:", targetChan)
		if !ok {
			return
		}
		select {
		case data := <-targetChan.(chan map[string][]byte):
			fmt.Println("data", data)
			fmt.Println("c.ClientId:", c.ClientId)
			fmt.Println("data[\"targetId\"]:", string(data["targetId"]))
			if c.ClientId == string(data["targetId"]) { // try 发给指定用户
				c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				c.Conn.WriteMessage(websocket.TextMessage, data["message"])
				fmt.Println("success write")
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			c.Conn.WriteMessage(websocket.PingMessage, nil)
		}
	}
}
