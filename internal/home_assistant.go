package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type homeAssistant struct {
	url                            string
	token                          string
	availability_percentage_entity map[string]string
	status_entity                  map[string]string
	intervention_entity            map[string]string
}

const (
	entity_id            = "id"
	entity_friendly_name = "friendly_name"
)

func newHomeAssistant() homeAssistant {
	result := homeAssistant{}
	result.url = viper.GetString("ha_url")
	result.token = viper.GetString("ha_token")
	result.availability_percentage_entity = viper.GetStringMapString("ha_frontforce_availability_percentage_entity")
	result.status_entity = viper.GetStringMapString("ha_frontforce_status_entity")
	result.intervention_entity = viper.GetStringMapString("ha_frontforce_intervention_entity")

	return result
}

func (h homeAssistant) updateAvailabilityPercentageState(availabilityStat availabiltyStat) error {
	attr := map[string]interface{}{
		"editable":      true,
		"min":           0,
		"max":           100,
		"pattern":       "null",
		"mode":          "text",
		"icon":          "mdi:fire",
		"friendly_name": h.availability_percentage_entity[entity_friendly_name],
	}
	err := h.updateState(h.availability_percentage_entity, attr, fmt.Sprintf("%.2f", availabilityStat.Periods[0].PercentAvailable))
	if err != nil {
		return err
	}
	return nil
}

func (h homeAssistant) updateStatusState(currAvail currentAvailability) error {
	attr := map[string]interface{}{
		"editable":      true,
		"min":           0,
		"max":           100,
		"pattern":       "null",
		"mode":          "text",
		"icon":          "mdi:fire",
		"friendly_name": h.status_entity[entity_friendly_name],
		"color":         currAvail.Unavailability.UnavailabilityCode.Color,
		"text_color":    currAvail.Unavailability.UnavailabilityCode.TextColor,
	}
	err := h.updateState(h.status_entity, attr, currAvail.Unavailability.UnavailabilityCode.Name)
	if err != nil {
		return err
	}
	return nil
}

func (h homeAssistant) updateInterventionState(interv intervention) error {
	attr := map[string]interface{}{
		"editable":      true,
		"min":           0,
		"max":           100,
		"pattern":       "null",
		"mode":          "text",
		"icon":          "mdi:fire",
		"friendly_name": h.intervention_entity[entity_friendly_name],
	}
	err := h.updateState(h.intervention_entity, attr, interv.InterventionCode.Description)
	if err != nil {
		return err
	}
	return nil
}

func (h homeAssistant) updateState(entity map[string]string, attributes map[string]interface{}, value string) error {
	var bearer = "Bearer " + h.token

	postBody, err := json.Marshal(map[string]interface{}{
		"state":      value,
		"attributes": attributes,
	})
	if err != nil {
		log.Error().Err(err).Msg("home assistant - failed marshalling state update value")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s/%s", h.url, "api/states", entity[entity_id]), bytes.NewBuffer(postBody))
	if err != nil {
		log.Error().Err(err).Msg("home assistant - failed creating update state request")
		return err
	}

	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msgf("home assistant - failed updating %s state", entity[entity_id])
		return err
	}
	if resp.StatusCode != 200 {
		log.Error().Msgf("home assistant - expected code 200, received: %d", resp.StatusCode)
		return err
	}
	log.Info().Msgf("home assistant - successfully updated %s state", entity[entity_id])
	return nil
}
