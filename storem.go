package websocketplugin

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"log"
	"time"
)

/*
 * storeInMongo
 * 消息存储到MongoDB中
 * 每道消息保存{ key(clientId-targetId; targetId-clientId), content(消息相关内容), time(当前时间) }
 * @param clientId 客户端ID; targetId 发送目标客户端ID; msgDetail 需要存储的消息相关内容
 */
func storeInMongo(clientId, targetId, msgDetail string) error {
	mcoll := GetMCOLL()
	fmt.Println("mcoll:", *mcoll)

	key1 := fmt.Sprintf("%s-%s", clientId, targetId) // clientId-targetId
	key2 := fmt.Sprintf("%s-%s", targetId, clientId) // targetId-clientId
	timeStamp := time.Now().Unix()

	result, err := mcoll.InsertOne(context.TODO(), bson.D{
		{"key", key1},
		{"content", msgDetail},
		{"timeStamp", timeStamp},
	})
	fmt.Println("result:", result)
	if err != nil {
		log.Println("存储到MongoDB失败: "+clientId+"to"+targetId, "\terror is: ", err)
		return errors.New("存储到MongoDB失败: " + clientId + "to" + targetId)
	}

	_, err = mcoll.InsertOne(context.TODO(), bson.D{
		{"key", key2},
		{"content", msgDetail},
		{"timeStamp", timeStamp},
	})
	if err != nil {
		log.Println("存储到MongoDB失败: "+targetId+"to"+clientId, "\terror is: ", err)
		return errors.New("存储到MongoDB失败: " + targetId + "to" + clientId)
	}

	return nil
}
