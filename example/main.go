package main

import (
	// "github.com/Ating56/websocket-plugin"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
	wp "websocket-plugin"
)

var Redis = wp.Redis{
	Host:     "127.0.0.1:6379",
	Password: "",
	DB:       0,
}
var MongoDB = wp.MongoDB{
	DataSource: "mongodb://root:student@127.0.0.1:27017",
	DataBase:   "test",
	Collection: "websocket",
}

type Target struct {
	TargetId string
	Content  string
}

func main() {
	wp.NewConfig(Redis, MongoDB)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		clientId := r.Header.Get("Sec-WebSocket-Protocol")
		wp.SetConnect(w, r, clientId)
	})
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		var target Target
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		defer r.Body.Close()
		err = json.Unmarshal(body, &target)
		if err != nil {
			return
		}
		fmt.Println("body", target.TargetId)

		wp.SendToServer(target.TargetId, target.Content)
	})
	http.Handle("/", http.FileServer(http.Dir("client")))
	http.ListenAndServe(":8080", nil)
}
