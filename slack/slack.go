package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type oAuthToken struct {
	Ok    bool
	Token string `json:"access_token"`
	Id    string `json:"user_id"`
	Team  struct {
		Id   string
		Name string
	}
	Authed_user struct {
		Id           string
		Access_token string
	}
}

type User struct {
	Id      string
	IsAdmin bool `json:"is_admin"`
	Deleted bool
	TeamId  string `json:"team_id"`

	Profile struct {
		Email string
	}
}

type userList struct {
	Ok       bool
	Error    string
	Needed   string
	Provided string
	Members  []User
}

func ExchangeToken(token string) (data oAuthToken) {
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

func UsersList(access_token string) (data []User, error error) {
	var d userList
	url := url.URL{
		Scheme: "https",
		Host:   "slack.com",
		Path:   "api/users.list",
	}

	client := &http.Client{}

	request, err := http.NewRequest("GET", url.String(), nil)

	if err != nil {
		return data, err
	}

	request.Header.Set("authorization", "Bearer "+access_token)

	res, err := client.Do(request)

	if err != nil {
		return data, err
	}

	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&d)

	if !d.Ok {
		fmt.Println("Error: ", d.Error)
		fmt.Println("Needed: ", d.Needed)
		fmt.Println("Provided: ", d.Provided)
		return data, err
	}

	if err != nil {
		return data, err
	}

	return d.Members, nil
}
