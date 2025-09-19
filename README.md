# websocket-plugin
Websocket plugin in Go
You can use it in your program to implement websocket chat.

- 基于gorilla/websocket
- upgrader设置了 Subprotocols: []string{r.Header.Get("Sec-WebSocket-Protocol")}，用于身份认证
- 前端建立连接时ws = new WebSocket('ws://127.0.0.1:8080/ws', ['1'])，传递第二个参数，可传递token，使用r.Header.Get("Sec-WebSocket-Protocol")接收
