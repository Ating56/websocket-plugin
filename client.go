package websocketplugin

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ClientId   string // 等效于用户Id
	Conn       *websocket.Conn
	Send       chan map[string][]byte
	RemoteAddr string
}

var GlobalClient *Client

// ReadPump 持续监听并处理客户端消息
// 实现原理：
// 1. 设置心跳超时（60秒）和Pong响应处理器
// 2. 读取消息并写入消息通道
// 3. 异常时注销连接并关闭通道
// 注意：
// - 使用defer确保连接最终关闭
// - 消息格式化为JSON包含客户端IP
func (c *Client) ReadPump() {
	defer func() {
		GlobalHub.unregister <- GlobalClient
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		msg := fmt.Sprintf("{\"ip\":\"%s\",\"message\":\"%s\"}", c.RemoteAddr, message)
		msgMap := map[string][]byte{
			"clientId": []byte(c.ClientId),
			"message":  []byte(msg),
		}
		select {
		case c.Send <- msgMap:
		default:
			log.Printf("Client %s queue full, disconnecting", c.RemoteAddr)
			close(c.Send)
		}
		log.Printf("Received from %s: %s", c.RemoteAddr, message)
	}
}

// WritePump 持续发送消息到客户端
// 工作流程：
// 1. 创建30秒间隔的心跳ticker
// 2. 监听消息通道和心跳事件
// 3. 处理消息发送或Ping心跳包
// 注意：
// - 设置10秒写超时防止阻塞
// - 通道关闭时退出循环
func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case data := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			c.Conn.WriteMessage(websocket.TextMessage, data["message"])
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			c.Conn.WriteMessage(websocket.PingMessage, nil)
		}
	}
}
