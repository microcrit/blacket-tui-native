package user

import (
	"encoding/json"
	"io"
	"net/http"

	"crit.rip/blacket-tui/api"

	"crit.rip/blacket-tui/api/types/responses"
)

func ClaimReward(token string) responses.ClaimRewardResponse {
	const BASE = api.API_BASE

	client := api.GetClient()
	req, _ := http.NewRequest("GET", BASE+"/worker/claim", nil)
	req.Header.Set("Cookie", "token="+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	bodyResp, _ := io.ReadAll(resp.Body)

	doc := responses.ClaimRewardResponse{}
	err = json.Unmarshal(bodyResp, &doc)
	if err != nil {
		panic(err)
	}
	return doc
}
