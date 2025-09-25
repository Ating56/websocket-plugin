package websocketplugin

import (
	"encoding/json"
	"errors"
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
 * Message
 * SendToServer 参数的结构体
 * 序列化后传入消息队列的数据信息
 */
type Message struct {

	// 客户端Id
	ClientId string

	// 目标客户端Id
	TargetId string

	// 消息相关的内容(发送者、接收者、消息内容、消息发送时间等信息)
	MsgDetail []byte

	// 消息发送时间戳
	TimeStamp int64
}

/*
 * GlobalRecv
 * 全局的接收通道, 对每个连接的client, 建立一个消息接收通道
 * clientId 作为key, 消息接收通道作为value
 */
var GlobalRecv sync.Map

/*
 * SendToServer
 * 客户端发送消息到服务端, 服务端检测接收方是否在线
 * 如果在线, 将消息写入接收方的消息接收通道并持久化存储
 * 如果不在线, 将消息持久化存储
 * @param clientId 连接的客户端Id; targetId 发送目标客户端Id; msgDetail 消息相关内容，[]byte格式
 */
func SendToServer(message Message) error {

	clientId, targetId, msgDetail, timeStamp := message.ClientId, message.TargetId, message.MsgDetail, message.TimeStamp

	if clientId == "" || targetId == "" {
		log.Println("SendToServer clientId or targetId is empty")
		return errors.New("clientId or targetId is empty")
	}

	// 二次组合消息内容, 加入targetId
	msg := map[string][]byte{
		"targetId":  []byte(targetId),
		"msgDetail": msgDetail,
	}
	// 目标客户端在线，写入目标客户端的消息接收通道
	if recvChan, ok := GlobalRecv.Load(targetId); ok {
		recvChan.(chan map[string][]byte) <- msg
	}

	// 存入redis
	err := storeInRedis(clientId, targetId, string(msgDetail))
	if err != nil {
		return err
	}

	// 存入mongo
	if GlobalMQInstance != nil {
		// 有MQ实例
		sendToMQmsg, marshalErr := json.Marshal(message)
		if marshalErr != nil {
			log.Println("json.Marshal error when send to MQ, err:", marshalErr)
			return marshalErr
		}
		// 发送到消息队列
		err = GlobalMQInstance.Publish(sendToMQmsg)
		if err != nil {
			log.Println("SendToServer Publish err:", err)
			return err
		}
	} else {
		// 无MQ实例
		err = storeInMongo(clientId, targetId, string(msgDetail), timeStamp)
		if err != nil {
			return err
		}
	}

	return nil
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
			log.Println("SendToClient's targetId:", string(data["targetId"]))
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
