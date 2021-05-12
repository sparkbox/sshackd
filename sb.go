package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sb/slack"
	"strings"
)

func filterUsers(users []slack.User) (filteredUsers []slack.User) {
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
			data := slack.ExchangeToken(code)
			users := slack.UsersList(data.Token)
			filteredUsers := filterUsers(users)
			json.NewEncoder(w).Encode(filteredUsers)
		}
	}

	http.HandleFunc("/", root)
	log.Fatal(http.ListenAndServe(":9001", nil))
}
