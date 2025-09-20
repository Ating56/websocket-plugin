package websocketplugin

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

/*
 * upgrader WebSocket升级器配置
 * 关键参数说明:
 * - Subprotocols: 子协议列表, 用于协商连接协议(当前配置为接收Sec-WebSocket-Protocol头)
 * - CheckOrigin: 跨域验证函数(当前配置允许所有来源，生产环境应限制域名)
 * 可扩展配置项:
 * - HandshakeTimeout: 握手超时时间(默认0-无限制)
 * - ReadBufferSize/WriteBufferSize: 读写缓冲区大小(单位字节)
 */
var upgrader = func(r *http.Request) *websocket.Upgrader {
	return &websocket.Upgrader{
		Subprotocols: []string{r.Header.Get("Sec-WebSocket-Protocol")},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

/*
 * SetConnect
 * 客户端连接到服务端
 * 升级HTTP连接为WebSocket连接
 * 注册新client到GlobalHub
 * @param w http.ResponseWriter; r * http.Request; clientId 连接的客户端ID
 */
func SetConnect(w http.ResponseWriter, r *http.Request, clientId string) {
	go GlobalHub.Run()
	conn, err := upgrader(r).Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	connectClient := &Client{
		ClientId:   clientId,
		Conn:       conn,
		RemoteAddr: r.RemoteAddr,
	}
	GlobalHub.register <- connectClient
}

/*
 * SetDisConnect
 * 客户端断开连接, 写入GlobalHub的unregister
 * @param clientId 连接的客户端ID
 */
func SetDisconnect(clientId string) {
	client, ok := GlobalHub.clients[clientId]
	if !ok {
		log.Printf("Client not found: %s\n", clientId)
		return
	}
	GlobalHub.unregister <- client
}
