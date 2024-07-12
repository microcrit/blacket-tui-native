package user

import (
	"encoding/json"
	"io"
	"net/http"

	"crit.rip/blacket-tui/api"

	"bytes"

	"crit.rip/blacket-tui/api/types/responses"
)

type SellRequest struct {
	Blook    string `json:"blook"`
	Quantity int    `json:"quantity"`
}

// Sell an item/blook immediately for tokens.
func Sell(token string, blookName string, quantity int) responses.GenericResponse {
	const BASE = api.API_BASE

	body := SellRequest{}
	body.Blook = blookName
	body.Quantity = quantity

	jsonValue, _ := json.Marshal(body)

	client := api.GetClient()
	req, _ := http.NewRequest("POST", BASE+"/worker/sell", bytes.NewBuffer(jsonValue))
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
