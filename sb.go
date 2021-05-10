package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type oAuthToken struct {
	Token string `json:"access_token"`
	Id    string `json:"user_id"`
}

type user struct {
	Id      string
	IsAdmin bool `json:"is_admin"`
	Deleted bool
}

type userList struct {
	Members []user
}

func exchangeToken(token string) (d oAuthToken) {
	var data oAuthToken

	url := url.URL{
		Scheme:   "https",
		Host:     "slack.com",
		Path:     "api/oauth.access",
		RawQuery: "client_id=" + os.Getenv("CLIENT_ID") + "&client_secret=" + os.Getenv("CLIENT_SECRET") + "&code=" + token + "&redirect_uri=https://sb-app-zvsm8hle.tunnelto.dev",
	}

	res, err := http.Get(url.String())

	if err != nil {
		log.Fatal("failed to request oauth token: ", err)
	}

	json.NewDecoder(res.Body).Decode(&data)

	return data
}

func get_info(access_token string) (data userList) {
	url := url.URL{
		Scheme: "https",
		Host:   "slack.com",
		Path:   "api/users.list",
	}

	client := &http.Client{}

	request, err := http.NewRequest("GET", url.String(), nil)

	if err != nil {
		log.Fatal("failed to create new request: ", err)
	}

	request.Header.Set("authorization", "Bearer "+access_token)

	res, err := client.Do(request)

	if err != nil {
		log.Fatal("failed to do request: ", err)
	}

	json.NewDecoder(res.Body).Decode(&data)

	defer res.Body.Close()

	if err != nil {
		log.Fatal("failed to parse user list: ", err)
	}

	return data
}

func main() {
	root := func(w http.ResponseWriter, req *http.Request) {
		tempCode := req.URL.Query()["code"]
		data := exchangeToken(strings.Join(tempCode, ""))
		_ = get_info(data.Token)
		io.WriteString(w, "transferring...\n")
	}

	http.HandleFunc("/", root)
	log.Fatal(http.ListenAndServe(":9001", nil))
}
