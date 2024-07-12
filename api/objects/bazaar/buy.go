package user

import (
	"encoding/json"
	"io"
	"net/http"

	"crit.rip/blacket-tui/api"

	"bytes"

	"crit.rip/blacket-tui/api/types/responses"
)

type BazaarBuyRequest struct {
	Id string `json:"id"`
}

// Buy an item/blook from the bazaar. An exact item id is required as returned by searching the bazaar.
func BazaarBuy(token string, itemId string) responses.GenericResponse {
	const BASE = api.API_BASE

	body := BazaarBuyRequest{}
	body.Id = itemId

	jsonValue, _ := json.Marshal(body)

	client := api.GetClient()
	req, _ := http.NewRequest("POST", BASE+"/worker/bazaar/buy", bytes.NewBuffer(jsonValue))
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
