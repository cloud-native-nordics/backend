package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	port    = 8080
	timeout = 2 * time.Second
)

func main() {
	slackToken := os.Getenv("SLACK_TOKEN")
	mux := http.NewServeMux()
	mux.HandleFunc("/invite", invite(slackToken))

	s := http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadTimeout:       timeout,
		WriteTimeout:      timeout,
		IdleTimeout:       timeout,
		ReadHeaderTimeout: timeout,
	}
	s.ListenAndServe()
}

func invite(slackToken string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		email := values.Get("email")

		if email == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		slackURL := "https://cloud-native-nordics.slack.com/api/users.admin.invite"

		values = url.Values{
			"email":      {email},
			"token":      {slackToken},
			"set_active": {"true"},
		}

		res, err := http.PostForm(slackURL, values)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&data)

		if err != nil {
			fmt.Printf("Decode: Something has gone wrong. Please contact a system administrator.\n")
			return
		}

		if data["ok"] == false {
			if data["error"].(string) == "already_invited" || data["error"].(string) == "already_in_team" {
				fmt.Printf("Success! You were already invited.\n")
				return
			} else if data["error"].(string) == "invalid_email" {
				fmt.Printf("The email you entered is an invalid email.\n")
				return
			} else if data["error"].(string) == "invalid_auth" {
				fmt.Printf("Invalid auth: Something has gone wrong. Please contact a system administrator.\n")
				return
			}
			fmt.Printf("Catch all: Something has gone wrong. Please contact a system administrator.\n")
			return
		}

		fmt.Printf("Success! Check “%s“ for an invite from Slack.\n", email)
		w.WriteHeader(http.StatusOK)
	}
}
