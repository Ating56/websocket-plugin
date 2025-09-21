package websocketplugin

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"time"
)

type MongoListRes struct {
	Id        bson.ObjectID `bson:"_id"`
	Key       string        `bson:"key"`
	Content   string        `bson:"content"`
	TimeStamp int64         `bson:"timeStamp"`
}

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

	// clientId与targetId一样，不重复存储
	if clientId == targetId {
		return nil
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

/*
 * GetMessageListInMongo
 * 获取Mongo中与目标客户端的消息列表
 * 分页, 倒序查询, lastMessageId传'', lastMessageTimeStamp传0, gap传0时, 查询最新10条记录
 * 否则查询_id<lastMessageId && timeStamp<lastMessageTimeStamp && 越过gap条的limit条数据的记录
 * @param clientId 客户端ID; targetId 目标客户端ID; lastMessageId 消息_id起点; lastMessageTimeStamp 消息时间戳起点; gap 越过的消息数量; limit 查询消息数量
 */
func GetMessageListInMongo(clientId, targetId, lastMessageId string, lastMessageTimeStamp, gap, limit int64) ([]MongoListRes, error) {
	mcoll := GetMCOLL()
	key := fmt.Sprintf("%s-%s", clientId, targetId)

	opts := options.Find().SetSort(bson.D{{"_id", -1}})
	if gap != 0 {
		opts.SetSkip(gap)
	}
	if limit != 0 {
		opts.SetLimit(limit)
	}

	// 筛选器
	filter, err := mongoFilterFunc(key, lastMessageId, lastMessageTimeStamp)
	if err != nil {
		return nil, err
	}

	res, err := mcoll.Find(context.TODO(), filter, opts)
	if err != nil {
		log.Println("查询MongoDB出错, error: ", err)
		return nil, errors.New("查询MongoDB出错")
	}

	var mongoListRes []MongoListRes
	if err = res.All(context.TODO(), &mongoListRes); err != nil {
		log.Println("查询MongoDB结果转化输出结构体出错, error: ", err)
		return nil, errors.New("查询MongoDB结果转化输出结构体出错")
	}

	return mongoListRes, nil
}

/*
 * mongoFilterFunc
 * mongo的筛选器, key精确查询, 小于lastMessageId, 小于lastMessageTimeStamp
 */
func mongoFilterFunc(key, lastMessageId string, lastMessageTimeStamp int64) (bson.D, error) {
	keyFilter, idFilter, timeStampFilter := bson.D{{"key", key}}, bson.D{}, bson.D{}

	// 传了lastMessageId, 加入lastMessageId筛选
	if lastMessageId != "" {
		idObjectID, err := bson.ObjectIDFromHex(lastMessageId)
		if err != nil {
			log.Println("lastMessageId转为bson.ObjectID失败, error: ", err)
			return bson.D{}, errors.New("lastMessageId转为bson.ObjectID失败")
		}
		idFilter = bson.D{{"_id", bson.D{{"$lt", idObjectID}}}}
	}

	// 传了lastMessageTimeStamp, 加入lastMessageTimeStamp筛选
	if lastMessageTimeStamp != 0 {
		timeStampFilter = bson.D{{"timeStamp", bson.D{{"$lt", lastMessageTimeStamp}}}}
	}

	return bson.D{
		{"$and",
			bson.A{
				keyFilter,
				idFilter,
				timeStampFilter,
			},
		},
	}, nil
}
