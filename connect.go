package websocketplugin

import (
	"net/http"
	"sync"
)

// 连接信息
type Config struct {
	Host     string // ip+port
	Route    string // route
	ClientId string
}

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

var once sync.Once
var (
	GlobalConfig  *Config
	GlobalRedis   *Redis
	GlobalMongoDB *MongoDB
)

func NewConnect(c ...any) (*Config, *Redis, *MongoDB) {
	once.Do(func() {
		for _, v := range c {
			switch v.(type) {
			case Config:
				GlobalConfig = &Config{
					Host:     v.(Config).Host,
					Route:    v.(Config).Route,
					ClientId: v.(Config).ClientId,
				}
				connect(GlobalConfig)
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
	return GlobalConfig, GlobalRedis, GlobalMongoDB
}

func connect(conf *Config) {
	go GlobalHub.Run()
	http.HandleFunc(conf.Route, WsHandler(GlobalHub, conf))
}
