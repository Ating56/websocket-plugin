package websocketplugin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
)

func TestSetConnect(t *testing.T) {
	clientId := "1"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := SetConnect(w, r, clientId)
		if err != nil {
			t.Errorf("SetConnect failed: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("connect success"))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	url := "ws" + server.URL[4:]
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Errorf("websocket Dial failed: %v", err)
	}
	defer ws.Close()
}
