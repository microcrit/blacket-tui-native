package data

import (
	"encoding/json"
	"io"
	"net/http"

	"crit.rip/blacket-tui/api"
)

func GetData() map[string]any {
	const BASE = api.API_BASE

	req, err := http.Get(BASE + "/data/index.json")
	if err != nil {
		panic(err)
	}

	defer req.Body.Close()

	body, _ := io.ReadAll(req.Body)

	doc := make(map[string]any)
	err = json.Unmarshal([]byte(body), &doc)

	if err != nil {
		panic(err)
	}

	return doc
}
