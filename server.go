package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	port    = 8080
	timeout = 2 * time.Second
)

func main() {
	slackToken := os.Getenv("SLACK_TOKEN")
	mux := http.NewServeMux()
	mux.HandleFunc("/invite", invite(slackToken))
	mux.HandleFunc("/ping", ping)

	s := http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadTimeout:       timeout,
		WriteTimeout:      timeout,
		IdleTimeout:       timeout,
		ReadHeaderTimeout: timeout,
	}
	log.Fatal(s.ListenAndServe())
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func invite(slackToken string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		email := values.Get("email")
		log.Infof("Email received: %s", email)

		if email == "" {
			fmt.Fprintf(w, "email not found")
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
			fmt.Fprintf(w, "slack post failed: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&data)

		if err != nil {
			fmt.Fprintf(w, "decode slack response failed")
			return
		}

		if data["ok"] == false {
			if data["error"].(string) == "already_invited" || data["error"].(string) == "already_in_team" {
				log.Infof("Success! You were already invited.\n")
				return
			} else if data["error"].(string) == "invalid_email" {
				log.Infof("The email you entered is an invalid email.\n")
				return
			} else if data["error"].(string) == "invalid_auth" {
				log.Infof("Invalid auth: Something has gone wrong. Please contact a system administrator.\n")
				return
			}
			log.Infof("Catch all: Something has gone wrong. Please contact a system administrator.\n")
			return
		}

		log.Infof("Success! Check “%s“ for an invite from Slack.\n", email)
		w.WriteHeader(http.StatusOK)
	}
}
