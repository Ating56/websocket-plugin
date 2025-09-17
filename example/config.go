package main

import (
	// "github.com/Ating56/websocket-plugin"
	"net/http"
	wp "websocket-plugin"
)

var RemoteAddr = wp.Config{
	Host:  "0.0.0.0:8080",
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

func main() {
	wp.NewConnect(RemoteAddr, Redis, MongoDB)
	http.ListenAndServe(":8080", nil)
}
