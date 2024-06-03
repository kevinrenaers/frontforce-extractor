package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	availabilityUrl  = "api/v1/person/getavailability"
	currentStatusUrl = "api/v1/unavailability/getcurrent"
	interventionUrl  = "api/v1/intervention/get"
)

type frontforce struct {
	auth            *authorization
	homeAssistant   homeAssistant
	url             string
	refreshInterval int
	fetchStats      bool
}

func NewFrontforce() (frontforce, error) {
	result := frontforce{
		auth:          newAuth(),
		homeAssistant: newHomeAssistant(),
		url:           viper.GetString("frontforce_url"),
		fetchStats:    viper.GetBool("fetch_stats"),
	}
	result.refreshInterval = viper.GetInt("refresh_interval")
	return result, nil
}

func (f frontforce) StartUpdater() {
	ticker := time.NewTicker(time.Duration(f.refreshInterval) * time.Second)
	for range ticker.C {
		log.Info().Msg("frontforce - updating home assistant values")
		f.updateHAValues()
	}
}

func (f frontforce) updateHAValues() {
	if f.fetchStats {
		availabilityStat, err := f.fetchAvailabilyStat()
		if err != nil {
			log.Error().Err(err).Msg("frontforce - failed fetching availability stat")
			return
		}
		err = f.homeAssistant.updateAvailabilityPercentageState(availabilityStat)
		if err != nil {
			log.Error().Err(err).Msg("frontforce - failed updating home assistant values")
		}
	}
	currAvail, err := f.fetchStatus()
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed fetching status")
		return
	}
	err = f.homeAssistant.updateStatusState(currAvail)
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed updating home assistant values")
	}
	currIntervention, err := f.fetchIntervention()
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed fetching intervention")
		return
	}
	err = f.homeAssistant.updateInterventionState(currIntervention)
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed updating home assistant values")
	}
}

func (f frontforce) fetchStatus() (currentAvailability, error) {
	var bearer = "Bearer " + f.auth.GetToken()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", f.url, currentStatusUrl), nil)
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

func (f frontforce) fetchAvailabilyStat() (availabiltyStat, error) {
	var bearer = "Bearer " + f.auth.GetToken()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", f.url, availabilityUrl), nil)
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed creating get availability statistics request")
		return availabiltyStat{}, err
	}
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed getting availability statistics")
		return availabiltyStat{}, err
	}
	if resp.StatusCode != 200 {
		log.Error().Msgf("frontforce - expected code 200, received: %d", resp.StatusCode)
		return availabiltyStat{}, err
	}
	decoder := json.NewDecoder(resp.Body)
	var availStatsResp availabiltyStat
	err = decoder.Decode(&availStatsResp)
	if err != nil {
		log.Error().Err(err).Msg("frontforce - failed decoding availability statistics")
		return availabiltyStat{}, err
	}
	log.Info().Msg("frontforce - successfully fetched frontforce availability statistics")
	return availStatsResp, nil
}

func (f frontforce) fetchIntervention() (intervention, error) {
	var bearer = "Bearer " + f.auth.GetToken()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", f.url, interventionUrl), nil)
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
