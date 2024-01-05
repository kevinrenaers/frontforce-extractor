package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type authorization struct {
	url         string
	username    string
	password    string
	jwtToken    string
	validUntil  time.Time
	refreshFrom time.Time
}

func newAuth() *authorization {
	auth := &authorization{}
	auth.url = viper.GetString("frontforce_url")
	auth.username = viper.GetString("frontforce_username")
	auth.password = viper.GetString("frontforce_password")
	auth.fetchToken()
	auth.startTicker()
	return auth
}

func (a *authorization) fetchToken() {
	now := time.Now()
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", a.username)
	data.Set("password", a.password)

	client := &http.Client{}
	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", a.url, "token"), strings.NewReader(data.Encode()))
	if err != nil {
		log.Error().Err(err).Msg("authorization - failed creating login request")
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	if err != nil {
		log.Error().Err(err).Msg("authorization - failed logging in")
	}
	if resp.StatusCode != 200 {
		log.Error().Msgf("authorization - expected code 200, received: %d", resp.StatusCode)
	}
	decoder := json.NewDecoder(resp.Body)
	var lr loginResp
	err = decoder.Decode(&lr)
	if err != nil {
		log.Error().Err(err).Msg("authorization - failed decoding login response")
	}

	a.jwtToken = lr.AccessToken
	a.validUntil = now.Add(time.Second * time.Duration(lr.ExpiresIn))
	a.refreshFrom = a.validUntil.Add(time.Minute * -15)
	log.Info().Msgf("authorization - token valid until: %s", a.validUntil.Format(time.RFC3339))
}

func (a *authorization) startTicker() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			now := time.Now()
			if !now.Before(a.refreshFrom) {
				log.Debug().Msg("authorization - token expired, fetching new one")
				a.fetchToken()
			}
		}
	}()
}

func (a authorization) GetToken() string {
	return a.jwtToken
}

type loginResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}
