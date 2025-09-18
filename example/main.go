package main

import (
	// "github.com/Ating56/websocket-plugin"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	wp "websocket-plugin"
)

var once sync.Once

var RemoteAddr = wp.Config{
	Route: "/ws",
}
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

type Client struct {
	ClientId string
}

func main() {
	wp.NewConfig(Redis, MongoDB)
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		defer r.Body.Close()
		var client Client
		err = json.Unmarshal(body, &client)
		if err != nil {
			return
		}
		fmt.Println("body", string(client.ClientId))
		once.Do(func() { wp.NewConnect(&RemoteAddr, client.ClientId) })
	})
	http.Handle("/", http.FileServer(http.Dir("client")))
	http.ListenAndServe(":8080", nil)
}
