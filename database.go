package websocketplugin

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	RDB   *redis.Client
	MDB   *mongo.Client
	MCOLL *mongo.Collection
)

// 连接Redis
func InitRDB(conf *Redis) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     conf.Host,
		Password: conf.Password,
		DB:       conf.DB,
	})
	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		panic("failed to connect redis: " + err.Error())
	}
	log.Println("success to connect redis")
}

// 连接MongoDB
func InitMDB(conf *MongoDB) {
	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(conf.DataSource).SetServerAPIOptions(serverApi)
	mdb, err := mongo.Connect(opts)
	if err != nil {
		panic("failed to connect MongoDB: " + err.Error())
	}
	defer func() {
		if err := mdb.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	MDB = mdb
	MCOLL = mdb.Database(conf.DataBase).Collection(conf.Collection)
	log.Println("success to connect MongoDB")
}

// 获取全局Redis
func GetRDB() *redis.Client {
	return RDB
}

// 获取全局MongoDB
func GetMDB() *mongo.Client {
	return MDB
}

// 获取全局MongoDB Collection
func GetMCOLL() *mongo.Collection {
	return MCOLL
}
