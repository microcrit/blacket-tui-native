package user

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"crit.rip/blacket-tui/api"
	"crit.rip/blacket-tui/util"

	"crit.rip/blacket-tui/api/types/objects"
	"crit.rip/blacket-tui/api/types/responses"

	"github.com/gbin/goncurses"
)

func Login(username string, password string) string {
	values := map[string]string{
		"username": username,
		"password": password,
	}
	jsonValue, _ := json.Marshal(values)

	const BASE = api.API_BASE

	resp, err := http.Post(BASE+"/worker/login", "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	cookie := resp.Header.Get("Set-Cookie")

	doc := responses.GenericResponse{}
	err = json.Unmarshal([]byte(body), &doc)
	if err != nil {
		panic(err)
	}

	if doc.Error {
		panic(doc.Reason)
	}

	return util.ParseCookie(cookie)
}

func LoginProxy(username string, password string, proxy string) string {
	values := map[string]string{
		"username": username,
		"password": password,
	}
	jsonValue, _ := json.Marshal(values)

	const BASE = api.API_BASE

	client := &http.Client{}
	url, err := url.Parse(proxy)
	if err != nil {
		panic(err)
	}
	client.Transport = &http.Transport{
		Proxy: http.ProxyURL(url),
	}

	resp, err := client.Post(BASE+"/worker/login", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	doc := responses.GenericResponse{}
	err = json.Unmarshal([]byte(body), &doc)
	if err != nil {
		panic(err)
	}

	if doc.Error {
		panic(doc.Reason)
	}

	return util.ParseCookie(resp.Header.Get("Set-Cookie"))
}

func GetUserBase(stdscr *goncurses.Window, token string, path string) responses.UserResponse {
	const BASE = api.API_BASE

	req, err := http.NewRequest("GET", BASE+path, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Cookie", "token="+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	doc := responses.UserResponse{}
	err = json.Unmarshal([]byte(body), &doc)
	if err != nil {
		panic(err)
	}

	return doc
}

// Get the user object of the currently authenticated user.
func GetUser(stdscr *goncurses.Window, token string) objects.User {
	return GetUserBase(stdscr, token, "/worker2/user").User
}

// Get the user object of a user by their username or ID.
func GetExternalUser(stdscr *goncurses.Window, token string, user util.Either[string, int]) objects.User {
	return util.Switch(user,
		func(username string) objects.User {
			return GetUserBase(stdscr, token, "/worker2/user/"+username).User
		},
		func(id int) objects.User {
			return GetUserBase(stdscr, token, "/worker2/user/"+strconv.Itoa(id)).User
		},
	)
}
