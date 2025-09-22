package websocketplugin

import (
	"sync"
)

// Redis
type Redis struct {
	Host     string // ip+port
	Password string // pwd
	DB       int    // 几号数据库
}

// MongoDB
type MongoDB struct {
	DataSource string // 数据源, 格式: mongodb://root:your_password@127.0.0.1:27017
	DataBase   string // 数据库名
	Collection string // 集合名
}

var onceConfig sync.Once
var (
	GlobalRedis   *Redis
	GlobalMongoDB *MongoDB
)

/*
 * NewConfig
 * 初始化配置
 * 接收Redis和MongoDB配置参数
 * 初始化Redis和MongoDB连接
 * @param Redis配置参数结构体; MongoDB配置参数结构体
 */
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
