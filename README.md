# websocket-plugin
基于gorilla/websocket二次开发的一个Go语言的websocket方法库，只需调用方法即可实现websocket相关功能。

## 功能
- [x] 建立websocket连接
- [x] 给目标客户端发送消息
- [x] 已连接的客户端自动接收消息
- [x] 查看消息记录

## 使用
### 添加配置
1. 添加Redis和MongoDB配置，调用NewConfig，参数为Redis和MongoDB配置，结构体如下：
    ```go
    // Redis配置结构体
    type Redis struct {
        Host     string // ip+port
        Password string // pwd
        DB       int    // 几号数据库
    }

    // MongoDB配置结构体
    type MongoDB struct {
        DataSource string // 数据源, 格式: mongodb://root:your_password@127.0.0.1:27017
        DataBase   string // 数据库名
        Collection string // 集合名
    }
    ```
2. 示例：
    ```go
    var Redis = wp.Redis{
        Host:     "127.0.0.1:6379",
        Password: "",
        DB:       0,
    }
    var MongoDB = wp.MongoDB{
        DataSource: "mongodb://root:your_password@127.0.0.1:27017",
        DataBase:   "test",
        Collection: "websocket",
    }

    wp.NewConfig(Redis, MongoDB)
    ```

### 建立连接
1. web-socket-plugin的upgrader设置了Subprotocols: []string{r.Header.Get("Sec-WebSocket-Protocol")}，用于用户身份认证
2. 客户端建立连接并传递token参数
3. 服务端接收参数并获取用户信息，并调用SetConnect，参数为：
    - http.ResponseWriter
    - *http.Request
    - clientId(客户端Id)
4. 已建立连接的客户端，实时接收消息
5. 示例：
    ```js
    // 客户端建立连接并传递token参数
    ws = new WebSocket('ws://127.0.0.1:8080/ws', [token])
    ```

    ```go
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        // 服务端接收参数
        token := r.Header.Get("Sec-WebSocket-Protocol")
        
        // ... 获取唯一的用户id，作为第三个参数传给SetConnect

		wp.SetConnect(w, r, clientId)
	})
    ```

### 发送消息
1. 发送消息是http请求，调用SendToServer，参数为：
    - clientId(客户端Id)
    - targetId(目标客户端Id)
    - msgDetail(消息相关信息，[]byte类型)
2. 示例：
    ```go
    data := `{"senderId": "1", "ReceiveId": "2", "content": "hello"}`
    err := SendToServer(clientId, targetId, []byte(data))
    ```

### 查询消息记录Redis
1. 消息存入Redis的数据格式为list，key的值为'clientId-targetId'，list的每个元素为消息的json字符串
2. 消息发送成功会存入Redis，每个list最多留存10条消息
3. 可调用GetMessageListInRedis获取Redis中的消息记录，参数为：
    - clientId(客户端Id)
    - targetId(目标客户端Id)
4. 示例：
    ```go
    res := GetMessageListInRedis(clientId, targetId)
    ```

### 查询消息记录MongoDB
1. 消息存入MongoDB的数据格式为collection中的一个文档，文档的字段有_id，key(格式为'clientId-targetId')，content(消息的json字符串)，timeStamp(发送消息的时间戳)
2. 消息发送成功会同时存入MongoDB，进行持久化存储
3. 可调用GetMessageListInMongo获取MongoDB中的消息记录(倒序查找)，参数为：
    - clientId(客户端Id)
    - targetId(目标客户端Id)
    - lastMessageId(最后一条消息Id，查找小于此Id的数据，若传空字符串则从最新一条开始查找)
    - lastMessageTimeStamp(最后一条消息的时间戳，查找小于此时间戳的数据，若传0则从最新一条开始查找)
    - gap(越过的消息条数，一般使用lastMessageId和lastMessageTimeStamp就可以实现越过一些消息，使用gap可保证与Redis查找的消息去重)
    - limit(查找的消息条数)
4. 示例：
    ```go
    res, err := GetMessageListInMongo(clientId, targetId, "", 0, 0, 10)
    ```
