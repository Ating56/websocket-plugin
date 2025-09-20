package websocketplugin

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

/*
 * storeInRedis
 * 将消息持久化到redis, 便于快速加载
 * 但redis中仅保存部分消息(暂定10条)
 * 每个聊天通道保存为一个list, 每条消息添加到两份list中(clientId-targetId; targetId-clientId)
 * @param clientId 客户端ID; targetId 发送目标客户端ID; msgDetail 需要存储的消息相关内容
 */
func storeInRedis(clientId, targetId, msgDetail string) error {
	rdb := GetRDB()
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
		return errors.New("存储Redis失败: " + clientId + "to" + targetId)
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
		return errors.New("存储Redis失败: " + targetId + "to" + clientId)
	}

	return nil
}
