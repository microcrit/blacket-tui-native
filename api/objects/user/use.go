package user

import (
	"encoding/json"
	"io"
	"net/http"

	"crit.rip/blacket-tui/api"

	"bytes"

	"crit.rip/blacket-tui/api/types/responses"
)

type UseRequest struct {
	Item string `json:"item"`
}

func Use(token string, item string) responses.GenericMessageResponse {
	const BASE = api.API_BASE

	body := UseRequest{}
	body.Item = item

	jsonValue, _ := json.Marshal(body)

	client := api.GetClient()
	req, _ := http.NewRequest("POST", BASE+"/worker/use", bytes.NewBuffer(jsonValue))
	req.Header.Set("Cookie", "token="+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	bodyResp, _ := io.ReadAll(resp.Body)

	doc := responses.GenericMessageResponse{}
	err = json.Unmarshal(bodyResp, &doc)
	if err != nil {
		panic(err)
	}
	return doc
}
