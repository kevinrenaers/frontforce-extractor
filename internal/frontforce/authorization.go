package frontforce

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"renaers.be/frontforce/internal/configuration"
)

type authorization struct {
	url          string
	refreshToken string
	tokenChannel chan string

	validUntil  time.Time
	refreshFrom time.Time
}

func newAuth(config configuration.Configuration) *authorization {
	auth := &authorization{
		url:          config.FrontforceURL,
		refreshToken: config.FrontforceInitialRefreshToken,
		tokenChannel: make(chan string),
	}

	return auth
}

func (a *authorization) Start() {
	a.fetchToken()
	a.startTicker()
}

func (a *authorization) TokenChannel() chan string {
	return a.tokenChannel
}

func (a *authorization) fetchToken() {
	now := time.Now()
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", a.refreshToken)
	data.Set("client_id", "EmergencyMobile")
	data.Set("device", `{"token":"pRLak92DxVLwN9HGiWXGqJC64AMv0UDUbywiWas2JHjNtCaj7I","language":"nl-BE","appType":1}`)

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

		return
	}
	decoder := json.NewDecoder(resp.Body)
	var lr loginResp
	err = decoder.Decode(&lr)
	if err != nil {
		log.Error().Err(err).Msg("authorization - failed decoding login response")
	}

	a.tokenChannel <- lr.AccessToken
	a.validUntil = now.Add(time.Second * time.Duration(lr.ExpiresIn))
	a.refreshToken = lr.RefreshToken
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

type loginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}
