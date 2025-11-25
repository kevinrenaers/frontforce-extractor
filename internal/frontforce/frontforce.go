package frontforce

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"renaers.be/frontforce/internal/configuration"
	"renaers.be/frontforce/internal/homeassistant"
	"renaers.be/frontforce/internal/influx"
	"renaers.be/frontforce/internal/shared"
	"renaers.be/frontforce/internal/websocket"
)

const (
	negotiateUrl     = "hub/dashboard/negotiate?dashboardId=%s&negotiateVersion=%s"
	dashboardUrl     = "hub/dashboard?dashboardId=%s"
	availabilityUrl  = "api/v1/person/getavailability"
	currentStatusUrl = "api/v1/unavailability/getcurrent"
	interventionUrl  = "api/v1/intervention/get"
	vehicleStatusUrl = "api/v1/unit/get"
)

type negotiateResponse struct {
	ConnectionToken string `json:"connectionToken"`
}

type frontforce struct {
	auth            *authorization
	homeAssistant   homeassistant.HomeAssistant
	influx          influx.Influx
	websocketClient websocket.WebsocketClient
	tokenChannel    chan string

	url    string
	userID int

	state *shared.State
}

func NewFrontforce(config configuration.Configuration) (frontforce, error) {
	dashboardUrl := fmt.Sprintf(dashboardUrl, config.FrontforceDashboardID)
	strippedBaseUrl := strings.TrimPrefix(config.FrontforceURL, "https://")

	result := frontforce{
		auth:            newAuth(config),
		homeAssistant:   homeassistant.NewHomeAssistant(config),
		influx:          influx.NewInflux(config),
		websocketClient: websocket.NewWebSocketClient(fmt.Sprintf("wss://%s/%s", strippedBaseUrl, dashboardUrl)),

		url:    config.FrontforceURL,
		userID: config.FrontforceUserID,

		state: shared.NewState(config),
	}

	go result.handleWebsocketMessages()

	return result, nil
}

func (f frontforce) Start() {
	go f.auth.Start()

	for token := range f.auth.TokenChannel() {
		err := f.websocketClient.Close()
		if err != nil {
			log.Error().Err(err).Msg("frontforce - failed closing websocket client")
		}

		dashboardToken, err := f.negotiate(token)
		if err != nil {
			log.Error().Err(err).Msg("frontforce - failed connecting websocket client")

			continue
		}

		err = f.websocketClient.Connect(dashboardToken, token)
		if err != nil {
			log.Error().Err(err).Msg("frontforce - failed connecting websocket client")
		}
	}
}

func (f frontforce) handleWebsocketMessages() {
	for msg := range f.websocketClient.Messages() {
		frontforceMessage := frontforceMessage{}

		err := json.Unmarshal(msg, &frontforceMessage)
		if err != nil {
			log.Error().Err(err).Msg("frontforce - failed unmarshalling websocket message")

			continue
		}

		relevantBlocks := frontforceMessage.FetchRelevantBlocks()

		if len(relevantBlocks) == 0 {
			log.Info().Msgf("frontforce - no relevant blocks found in websocket message, so skipping updates")

			continue
		}

		log.Info().Msgf("frontforce - fetched %d relevant blocks from websocket message", len(relevantBlocks))

		err = f.state.Update(relevantBlocks)
		if err != nil {
			log.Error().Err(err).Msg("frontforce - failed updating frontforce state")

			continue
		}

		f.updateHAValues(*f.state)
	}
}

func (f frontforce) updateHAValues(state shared.State) {
	person := state.GetPerson(f.userID)
	err := f.homeAssistant.UpdateStatusState(person.UnavailabilityCode)
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed updating home assistant personal status values")
	}
	interventions := state.Interventions
	err = f.homeAssistant.UpdateInterventionState(interventions)
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed updating home assistant interventions values")
	}
	err = f.homeAssistant.UpdateVehicleStatsState(state.Vehicles)
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed updating home assistant vehicles values")
	}
}

func (f frontforce) negotiate(accessToken string) (string, error) {
	resp := negotiateResponse{}
	var bearer = "Bearer " + accessToken

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", f.url, fmt.Sprintf(negotiateUrl, "61", "1")), nil)
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed creating dashboard negotiation request")
		return "", err
	}
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	httpResp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed negotiating dashboard")

		return "", err
	}
	if httpResp.StatusCode != 200 {
		log.Error().Msgf("expected code 200, received: %d", httpResp.StatusCode)

		return "", err
	}
	decoder := json.NewDecoder(httpResp.Body)
	err = decoder.Decode(&resp)
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed decoding dashboard negotiation response")
		return "", err
	}

	log.Info().Msg("frontforce - successfully negotiated dashboard")

	return resp.ConnectionToken, nil
}
