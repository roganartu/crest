package crest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"

	"golang.org/x/oauth2"
)

const (
	CALLBACK_URL = "http://localhost:8080/callback"
	AUTH_BASE    = "https://login.eveonline.com"
)

var (
	DEFAULT_SCOPE       = []string{"publicData"}
	callback            net.Listener
	callbackStatus      chan string
	callbackErrorRegexp = regexp.MustCompile(`^error:`)
	client              *http.Client
)

// Authenticate calls the oauth endpoint and spins up a HTTP server to handle the callback.
func Authenticate(clientID, clientSecret string, scopes []string) (*http.Client, error) {
	if len(scopes) < 1 {
		scopes = DEFAULT_SCOPE
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  CALLBACK_URL,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  AUTH_BASE + "/oauth/authorize",
			TokenURL: AUTH_BASE + "/oauth/token",
		},
	}

	// Try and load token from cache
	f, err := os.Open(".token")
	if err == nil {
		var t oauth2.Token
		err = json.NewDecoder(f).Decode(&t)
		f.Close()
		if err == nil {
			client = config.Client(oauth2.NoContext, &t)
			return client, nil
		}
	}

	// Spawn a listener that we can close, killing the http server we spawn.
	callback, err = net.Listen("tcp", ":8080")
	if err != nil {
		return nil, err
	}

	callbackStatus = make(chan string)
	defer close(callbackStatus)

	fmt.Println(config.AuthCodeURL("state", oauth2.AccessTypeOffline))

	http.HandleFunc("/callback", callbackHandler)
	http.Serve(callback, nil)

	status := <-callbackStatus
	if callbackErrorRegexp.MatchString(status) {
		return nil, errors.New(status)
	}

	token, err := config.Exchange(oauth2.NoContext, status)
	if err != nil {
		return nil, err
	}

	// Cache the token for later use
	f, err = os.Create(".token")
	if err == nil {
		err = json.NewEncoder(f).Encode(token)
		f.Close()
	}

	client = config.Client(oauth2.NoContext, token)
	return client, nil
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	callback.Close()
	if r.FormValue("state") != "state" {
		callbackStatus <- "Invalid state verification string."
		return
	}

	callbackStatus <- r.FormValue("code")
}
