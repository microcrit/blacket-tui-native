package user

import (
	"encoding/json"
	"io"
	"net/http"

	"crit.rip/blacket-tui/api"

	"bytes"

	"crit.rip/blacket-tui/api/types/responses"
)

type BazaarListRequest struct {
	Item  string `json:"item"`
	Price int    `json:"price"`
}

// List an item/blook on the bazaar for a given price.
func BazaarList(token string, item string, price int) responses.GenericResponse {
	const BASE = api.API_BASE

	body := BazaarListRequest{}
	body.Item = item
	body.Price = price

	jsonValue, _ := json.Marshal(body)

	client := api.GetClient()
	req, _ := http.NewRequest("POST", BASE+"/worker/bazaar/list", bytes.NewBuffer(jsonValue))
	req.Header.Set("Cookie", "token="+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	bodyResp, _ := io.ReadAll(resp.Body)

	doc := responses.GenericResponse{}
	err = json.Unmarshal(bodyResp, &doc)
	if err != nil {
		panic(err)
	}
	return doc
}
