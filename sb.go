package main

import (
	"encoding/json"
	"io"
	"net/http"
	"sb/ca"
	"sb/slack"
	"strings"
)

type SSHLogin struct {
	Token string
}

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

func isValidUser(users []slack.User, userId string) (found bool) {
	for i := 0; i < len(users); i++ {
		id := users[i].Id
		if id == userId {
			return true
		}
	}
	return false
}

func login(w http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")

	if code != "" {
		data := slack.ExchangeToken(code)
		io.WriteString(w, data.Authed_user.Id+":"+data.Authed_user.Access_token)
	}
}

func ssh(w http.ResponseWriter, req *http.Request) {
	var s SSHLogin

	json.NewDecoder(req.Body).Decode(&s)

	defer req.Body.Close()

	splits := strings.Split(s.Token, ":")

	users := slack.UsersList(splits[1])
	filteredUsers := filterUsers(users)
	isValid := isValidUser(filteredUsers, splits[0])
	if isValid {
		key := ca.SignCert()
		io.WriteString(w, key)
	} else {
		w.WriteHeader(401)
	}
}

func main() {
	// http.HandleFunc("/login", login)
	// http.HandleFunc("/ssh", ssh)
	// log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
	_ = ca.SignCert()
}
