package websocketplugin

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ClientId   string // 等效于发送用户Id
	Conn       *websocket.Conn
	RemoteAddr string
}

/*
 * GlobalRecv
 * 全局的接收通道，对每个连接的client，建立一个消息接收通道
 * clientId 作为key，消息接收通道作为value
 */
var GlobalRecv sync.Map

/*
 * SendToServer
 * 客户端发送消息到服务端，服务端检测接收方是否在线
 * 如果在线，将消息写入接收方的消息接收通道并持久化存储
 * 如果不在线，将消息持久化存储
 */
func SendToServer(clientId, targetId string, msgDetail []byte) {
	msg := map[string][]byte{
		"targetId":  []byte(targetId),
		"msgDetail": msgDetail,
	}
	log.Println("SendToServer targetId:", targetId, "msgDetail:", msgDetail)
	if recvChan, ok := GlobalRecv.Load(targetId); ok {
		recvChan.(chan map[string][]byte) <- msg
	}
	// todo 持久化存储消息
}

/*
 * SendToClient
 * 服务端将消息发给客户端
 * 创建30秒间隔的心跳ticker
 * 监听消息通道和心跳事件
 * 处理消息发送或Ping心跳包
 * - 设置10秒写超时防止阻塞
 */
func (c *Client) SendToClient() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		targetChan, ok := GlobalRecv.Load(c.ClientId)
		if !ok {
			return
		}
		select {
		case data := <-targetChan.(chan map[string][]byte):
			log.Println("data[\"targetId\"]:", string(data["targetId"]))
			if c.ClientId == string(data["targetId"]) {
				c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				c.Conn.WriteMessage(websocket.TextMessage, data["msgDetail"])
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			c.Conn.WriteMessage(websocket.PingMessage, nil)
		}
	}
}
