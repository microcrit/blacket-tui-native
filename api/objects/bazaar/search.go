package user

import (
	"encoding/json"
	"io"
	"net/http"

	"crit.rip/blacket-tui/api"

	"crit.rip/blacket-tui/api/types/responses"
)

// Search the bazaar for a specific item/blook or search by user ID. Providing a user ID will return listings from that user, otherwise the resulting items and blooks will contain the query in their name.
func BazaarSearch(token string, query string) responses.BazaarSearchResponse {
	const BASE = api.API_BASE

	client := api.GetClient()
	req, _ := http.NewRequest("GET", BASE+"/worker/bazaar/"+query, nil)
	req.Header.Set("Cookie", "token="+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	bodyResp, _ := io.ReadAll(resp.Body)

	doc := responses.BazaarSearchResponse{}
	err = json.Unmarshal(bodyResp, &doc)
	if err != nil {
		panic(err)
	}
	return doc
}
