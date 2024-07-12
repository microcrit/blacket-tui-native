package api

import (
	"net/http"
)

func GetClient() *http.Client {
	return &http.Client{}
}
