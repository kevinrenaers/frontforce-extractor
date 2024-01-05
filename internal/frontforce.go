package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type frontforce struct {
	auth            *authorization
	homeAssistant   homeAssistant
	url             string
	refreshInterval int
}

func NewFrontforce() (frontforce, error) {
	result := frontforce{
		auth:          newAuth(),
		homeAssistant: newHomeAssistant(),
		url:           viper.GetString("frontforce_url"),
	}
	result.refreshInterval = viper.GetInt("refresh_interval")
	return result, nil
}

func (f frontforce) StartUpdater() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		log.Info().Msg("frontforce - updating home assistant values")
		err := f.updateHAValues()
		if err != nil {
			panic(err)
		}
	}
}

func (f frontforce) updateHAValues() error {
	currAvail, err := f.fetchStatus()
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed fetching status")
		return err
	}
	err = f.homeAssistant.updateStatusState(currAvail)
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed updating home assistant values")
		return err
	}
	currIntervention, err := f.fetchIntervention()
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed fetching intervention")
		return err
	}
	err = f.homeAssistant.updateInterventionState(currIntervention)
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed updating home assistant values")
		return err
	}
	return nil
}

func (f frontforce) fetchStatus() (currentAvailability, error) {
	var bearer = "Bearer " + f.auth.GetToken()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", f.url, "api/v1/unavailability/getcurrent"), nil)
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed creating get current availability request")
		return currentAvailability{}, err
	}
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed getting current availability")
		return currentAvailability{}, err
	}
	if resp.StatusCode != 200 {
		log.Error().Msgf("frontforce - expected code 200, received: %d", resp.StatusCode)
		return currentAvailability{}, err
	}
	decoder := json.NewDecoder(resp.Body)
	var availResp currentAvailability
	err = decoder.Decode(&availResp)
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed decoding current availability")
		return currentAvailability{}, err
	}
	log.Info().Msg("frontforce - successfully fetched frontforce status")
	return availResp, nil
}

func (f frontforce) fetchIntervention() (intervention, error) {
	var bearer = "Bearer " + f.auth.GetToken()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", f.url, "api/v1/intervention/get"), nil)
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed creating get current intervention request")
		return intervention{}, err
	}
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed getting current intervention")
		return intervention{}, err
	}
	if resp.StatusCode != 200 {
		log.Error().Msgf("frontforce - expected code 200, received: %d", resp.StatusCode)
		return intervention{}, err
	}
	decoder := json.NewDecoder(resp.Body)
	var interv intervention
	err = decoder.Decode(&interv)
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed decoding intervention")
		return intervention{}, err
	}
	log.Info().Msg("frontforce - successfully fetched frontforce intervention")
	return interv, nil
}
