package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type oAuthToken struct {
	Ok    bool
	Token string `json:"access_token"`
	Id    string `json:"user_id"`
	Team  struct {
		Id   string
		Name string
	}
	authed_user struct {
		Id string
	}
}

type user struct {
	Id      string
	IsAdmin bool `json:"is_admin"`
	Deleted bool
	TeamId  string `json:"team_id"`

	Profile struct {
		Email string
	}
}

type userList struct {
	Members []user
}

func exchangeToken(token string) (data oAuthToken) {
	url := url.URL{
		Scheme:   "https",
		Host:     "slack.com",
		Path:     "api/oauth.v2.access",
		RawQuery: "client_id=" + os.Getenv("CLIENT_ID") + "&client_secret=" + os.Getenv("CLIENT_SECRET") + "&code=" + token,
	}

	res, err := http.Get(url.String())

	if err != nil {
		log.Fatal("request oauth token failed: ", err)
	}

	json.NewDecoder(res.Body).Decode(&data)

	if !data.Ok {
		e, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		log.Fatal(string(e))
	}

	defer res.Body.Close()

	return data
}

func users_list(access_token string) (data []user) {
	var d userList
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

	json.NewDecoder(res.Body).Decode(&d)

	defer res.Body.Close()

	if err != nil {
		log.Fatal("failed to parse user list: ", err)
	}

	return d.Members
}

func filterUsers(users []user) (filteredUsers []user) {
	for i := 0; i < len(users); i++ {
		emailDomain := returnDomain(users[i].Profile.Email)
		deleted := users[i].Deleted

		if emailDomain == "heysparkbox.com" && !deleted {
			filteredUsers = append(filteredUsers, users[i])
		}
	}

	return filteredUsers
}

// Adapter from net/mail
// https://golang.org/src/net/mail/message.go?s=5869:5932#L219
func returnDomain(address string) (domain string) {
	at := strings.LastIndex(address, "@")
	return address[at+1:]
}

func main() {
	root := func(w http.ResponseWriter, req *http.Request) {
		code := req.URL.Query().Get("code")

		if code != "" {
			data := exchangeToken(code)
			users := users_list(data.Token)
			filteredUsers := filterUsers(users)
			json.NewEncoder(w).Encode(filteredUsers)
		}
	}

	http.HandleFunc("/", root)
	log.Fatal(http.ListenAndServe(":9001", nil))
}
