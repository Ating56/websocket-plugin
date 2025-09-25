package websocketplugin

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

/*
 * storeInRedis
 * 将消息持久化到redis, 便于快速加载
 * 但redis中仅保存部分消息(暂定10条)
 * 每个聊天通道保存为一个list, 每条消息添加到两份list中(clientId-targetId; targetId-clientId)
 * @param clientId 客户端Id; targetId 发送目标客户端Id; msgDetail 需要存储的消息相关内容
 */
func storeInRedis(clientId, targetId, msgDetail string) error {
	rdb := GetRDB()
	if rdb == nil {
		log.Println("Redis未初始化")
		return errors.New("Redis未初始化")
	}

	key1 := fmt.Sprintf("%s-%s", clientId, targetId) // clientId-targetId
	key2 := fmt.Sprintf("%s-%s", targetId, clientId) // targetId-clientId

	_, err := rdb.TxPipelined(context.Background(), func(pipeliner redis.Pipeliner) error {
		len1 := pipeliner.LLen(context.Background(), key1).Val()
		if len1 > 10 {
			pipeliner.RPop(context.Background(), key1)
		}
		pipeliner.LPush(context.Background(), key1, msgDetail)
		return nil
	})
	if err != nil {
		log.Println("存储到Redis失败: "+clientId+"to"+targetId, "\terror is: ", err)
		return errors.New("存储到Redis失败: " + clientId + "to" + targetId)
	}

	// clientId与targetId一样，不重复存储
	if clientId == targetId {
		return nil
	}

	_, err = rdb.TxPipelined(context.Background(), func(pipeliner redis.Pipeliner) error {
		len2 := pipeliner.LLen(context.Background(), key2).Val()
		if len2 > 10 {
			pipeliner.RPop(context.Background(), key2)
		}
		pipeliner.LPush(context.Background(), key2, msgDetail)
		return nil
	})
	if err != nil {
		log.Println("存储到Redis失败: "+targetId+"to"+clientId, "\terror is: ", err)
		return errors.New("存储到Redis失败: " + targetId + "to" + clientId)
	}

	return nil
}

/*
 * GetMessageListInRedis
 * 获取redis中与目标客户端的消息列表
 * @param clientId 客户端Id; targetId 目标客户端Id
 */
func GetMessageListInRedis(clientId, targetId string) ([]string, error) {
	rdb := GetRDB()
	if rdb == nil {
		log.Println("Redis未初始化")
		return nil, errors.New("Redis未初始化")
	}

	key := fmt.Sprintf("%s-%s", clientId, targetId)

	res := rdb.LRange(context.Background(), key, 0, -1).Val()
	return res, nil
}
