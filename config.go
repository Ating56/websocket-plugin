package websocketplugin

import (
	"sync"
)

// Redis
type Redis struct {
	Host     string
	Password string
	DB       int
}

// MongoDB
type MongoDB struct {
	DataSource string
	DataBase   string
	Collection string
}

var onceConfig sync.Once
var (
	GlobalRedis   *Redis
	GlobalMongoDB *MongoDB
)

func NewConfig(c ...any) (*Redis, *MongoDB) {
	onceConfig.Do(func() {
		for _, v := range c {
			switch v.(type) {
			case Redis:
				GlobalRedis = &Redis{
					Host:     v.(Redis).Host,
					Password: v.(Redis).Password,
					DB:       v.(Redis).DB,
				}
				InitRDB(GlobalRedis)
			case MongoDB:
				GlobalMongoDB = &MongoDB{
					DataSource: v.(MongoDB).DataSource,
					DataBase:   v.(MongoDB).DataBase,
					Collection: v.(MongoDB).Collection,
				}
				InitMDB(GlobalMongoDB)
			}
		}
	})
	return GlobalRedis, GlobalMongoDB
}
