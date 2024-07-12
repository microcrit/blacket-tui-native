package chat

import (
	"log"

	"encoding/json"
	"net/http"

	"crit.rip/blacket-tui/api"
	"crit.rip/blacket-tui/api/types/objects"
	"github.com/gorilla/websocket"
)

var conn *websocket.Conn

func Handler(token string, user objects.User) {
	const BASE = api.API_BASE

	req, err := http.NewRequest("GET", BASE+"/worker/socket", nil)

	if err != nil {
		log.Fatal("req:", err)
		return
	}

	req.Header.Set("Cookie", "token="+token)

	conn, _, err := websocket.DefaultDialer.Dial(req.URL.String(), req.Header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()
}

func SendMessage(message string) {
	data := map[string]interface{}{
		"event": "messages-create",
		"data": map[string]interface{}{
			"room":    0,
			"content": message,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal("json:", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		log.Fatal("write:", err)
		return
	}
}
