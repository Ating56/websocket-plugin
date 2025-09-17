package websocketplugin

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocket升级器配置
// 关键参数说明：
// - CheckOrigin: 跨域验证函数（当前配置允许所有来源，生产环境应限制域名）
// 可扩展配置项：
// - HandshakeTimeout: 握手超时时间（默认0-无限制）
// - ReadBufferSize/WriteBufferSize: 读写缓冲区大小（单位字节）
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsHandler 处理WebSocket连接请求
// 流程说明：
// 1. 升级HTTP连接到WebSocket协议
// 2. 初始化客户端实例并注册到Hub
// 3. 启动读写协程维护双工通信
// 参数说明：
// - hub: 全局连接管理器实例
// - w: 响应写入器
// - r: 包含请求头等信息的HTTP请求对象
func WsHandler(hub *Hub, conf *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}

		GlobalClient = &Client{
			ClientId:   conf.ClientId,
			Conn:       conn,
			Send:       make(chan map[string][]byte, 256),
			RemoteAddr: r.RemoteAddr,
		}
		hub.register <- GlobalClient

		go GlobalClient.ReadPump()
		go GlobalClient.WritePump()
	}
}
