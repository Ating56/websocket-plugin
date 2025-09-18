package websocketplugin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// 连接信息
type Config struct {
	Route string // route
}
type Target struct {
	TargetId string
}

var onceConnect sync.Once

var GetClientId func() string

func NewConnect(conf *Config, f func() string) {
	go GlobalHub.Run()
	GetClientId = f
	onceConnect.Do(func() {
		http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("come send")
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return
			}
			defer r.Body.Close()
			var target Target
			err = json.Unmarshal(body, &target)
			if err != nil {
				return
			}
			targetId := target.TargetId
			fmt.Println("targetId:", targetId)
			sendToServer(targetId, "hello")
		})
		http.HandleFunc(conf.Route, WsHandler(GlobalHub, conf, f))
	})
}
