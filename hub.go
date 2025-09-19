package websocketplugin

import "log"

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

var GlobalHub *Hub = &Hub{
	clients:    make(map[*Client]bool),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

// Run 启动Hub的事件循环
// 采用select多路复用实现事件调度，处理流程：
// 1. 接收注册请求：将客户端加入映射表并记录日志
// 2. 接收注销请求：安全移除客户端并关闭消息通道
// 注意事项：
// - 使用无缓冲通道确保事件顺序处理
// - close(client.Send) 避免通道泄漏
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			log.Printf("New client registered: %s and clientId is: %s\n", client.RemoteAddr, client.ClientId)
			h.clients[client] = true
			GlobalRecv.Store(client.ClientId, make(chan map[string][]byte, 256))

			go client.SendToClient()
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				log.Printf("Client unregistered: %s\n", client.RemoteAddr)
				delete(h.clients, client)
				// if recvChan, ok := GlobalRecv.Load(client.ClientId); ok {
				// 	close(recvChan.(chan map[string][]byte)) // todo 优化 关闭通道持续触发 client.go的 data := <-targetChan.(chan map[string][]byte)
				// }
			}
		}
	}
}
