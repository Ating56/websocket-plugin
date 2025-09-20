# websocket-plugin
Websocket plugin in Go
You can use it in your program to implement websocket chat.

- 基于gorilla/websocket
- upgrader设置了 Subprotocols: []string{r.Header.Get("Sec-WebSocket-Protocol")}，用于身份认证
- 前端建立连接时ws = new WebSocket('ws://127.0.0.1:8080/ws', ['1'])，传递第二个参数，可传递token，使用r.Header.Get("Sec-WebSocket-Protocol")接收
- 暂不支持自我对话
- 最新一条消息存储到redis
- 所有消息存储到mongo

todo
- example接收到消息，判断是哪个对话的，添加未读/置顶；或弹出提示框XXX发来消息
- 消息持久化存储，发送消息需要把消息保存两份，一份存到自己对联系人的消息列表；另一份存到联系人对我的消息列表
    - 登录用户id{
        联系人id1: {
            内容
        },
    }
