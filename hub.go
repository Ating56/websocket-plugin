package websocketplugin

import "log"

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

/*
 * GlobalHub
 * 全局的Hub，管理所有连接的clients
 * client建立连接 -> 写入register
 * client关闭连接 -> 写入unregister
 */
var GlobalHub *Hub = &Hub{
	clients:    make(map[*Client]bool),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

// Run 启动服务就开启协程，GlobalHub.Run()，监听client的注册和注销
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
