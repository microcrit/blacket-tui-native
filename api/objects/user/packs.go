package user

import (
	"encoding/json"
	"io"
	"net/http"

	"crit.rip/blacket-tui/api"

	"bytes"

	"crit.rip/blacket-tui/api/types/responses"
)

type PackRequest struct {
	Pack string `json:"pack"`
}

// Open a pack from the market.
func OpenPack(token string, pack string) responses.PackOpenResponse {
	const BASE = api.API_BASE

	body := PackRequest{}
	body.Pack = pack

	jsonValue, _ := json.Marshal(body)

	client := api.GetClient()

	req, _ := http.NewRequest("POST", BASE+"/worker3/open", bytes.NewBuffer(jsonValue))
	req.Header.Set("Cookie", "token="+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	bodyResp, _ := io.ReadAll(resp.Body)

	doc := responses.PackOpenResponse{}
	err = json.Unmarshal([]byte(bodyResp), &doc)
	if err != nil {
		panic(err)
	}

	return doc
}
