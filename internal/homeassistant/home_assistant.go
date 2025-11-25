package homeassistant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"renaers.be/frontforce/internal/configuration"
	"renaers.be/frontforce/internal/shared"
)

type HomeAssistant struct {
	config configuration.Configuration
}

var (
	availabilityCodeIDBlackColor = []int{-1, 4, 7, 8, 9, 10, 13, 16, 17, 18, 19, 20, 23, 26, 27, 28, 32, 35, 37, 38, 41, 42, 43, 45, 52, 53, 54}
)

func NewHomeAssistant(config configuration.Configuration) HomeAssistant {
	result := HomeAssistant{
		config: config,
	}
	// result.url = config.HaUrl
	// result.token = config.HaToken
	// result.availability_percentage_entity = viper.GetStringMapString("ha_frontforce_availability_percentage_entity")
	// result.hours_needed_entity = viper.GetStringMapString("ha_frontforce_hours_needed_entity")
	// result.hours_performed_entity = viper.GetStringMapString("ha_frontforce_hours_performed_entity")
	// result.hours_remaining_entity = viper.GetStringMapString("ha_frontforce_hours_remaining_entity")
	// result.max_possible_availability_entity = viper.GetStringMapString("ha_frontforce_max_possible_availability_entity")
	// result.status_entity = viper.GetStringMapString("ha_frontforce_status_entity")
	// result.intervention_entity = viper.GetStringMapString("ha_frontforce_intervention_entity")
	// result.vehicle_status = viper.GetStringMapString("ha_frontforce_vehicle_status")

	return result
}

// func (h HomeAssistant) updateAvailabilityStates(availabilityStat internal.AvailabiltyStat) error {
// 	attr := map[string]interface{}{
// 		"editable":      false,
// 		"pattern":       "null",
// 		"mode":          "text",
// 		"icon":          "mdi:fire",
// 		"friendly_name": h.availability_percentage_entity[entity_friendly_name],
// 	}
// 	err := h.updateState(h.availability_percentage_entity[entity_id], attr, fmt.Sprintf("%.2f", availabilityStat.Periods[0].PercentAvailable))
// 	if err != nil {
// 		return err
// 	}

// 	hoursAvailable, neededHours, hoursToGo, maxPercentage := calcAvailabilityStats(availabilityStat.Periods[0].PercentAvailable)
// 	attr = map[string]interface{}{
// 		"editable":      false,
// 		"min":           0,
// 		"max":           100,
// 		"mode":          "box",
// 		"icon":          "mdi:fire",
// 		"step":          "0.001",
// 		"friendly_name": h.hours_performed_entity[entity_friendly_name],
// 	}
// 	err = h.updateState(h.hours_performed_entity[entity_id], attr, fmt.Sprintf("%.3f", hoursAvailable))
// 	if err != nil {
// 		return err
// 	}
// 	attr = map[string]interface{}{
// 		"editable":      false,
// 		"min":           0,
// 		"max":           100,
// 		"mode":          "box",
// 		"icon":          "mdi:fire",
// 		"step":          "0.001",
// 		"friendly_name": h.hours_needed_entity[entity_friendly_name],
// 	}
// 	err = h.updateState(h.hours_needed_entity[entity_id], attr, fmt.Sprintf("%.3f", neededHours))
// 	if err != nil {
// 		return err
// 	}
// 	attr = map[string]interface{}{
// 		"editable":      false,
// 		"min":           0,
// 		"max":           100,
// 		"mode":          "box",
// 		"icon":          "mdi:fire",
// 		"step":          "0.001",
// 		"friendly_name": h.hours_remaining_entity[entity_friendly_name],
// 	}
// 	err = h.updateState(h.hours_remaining_entity[entity_id], attr, fmt.Sprintf("%.3f", hoursToGo))
// 	if err != nil {
// 		return err
// 	}
// 	attr = map[string]interface{}{
// 		"editable":      false,
// 		"min":           0,
// 		"max":           100,
// 		"mode":          "box",
// 		"icon":          "mdi:fire",
// 		"step":          "0.001",
// 		"friendly_name": h.max_possible_availability_entity[entity_friendly_name],
// 	}
// 	err = h.updateState(h.max_possible_availability_entity[entity_id], attr, fmt.Sprintf("%.2f", maxPercentage))
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (h HomeAssistant) UpdateStatusState(currAvail shared.UnavailabilityCode) error {
	textColor := "#ffffff"
	if slices.Contains(availabilityCodeIDBlackColor, currAvail.ID) {
		textColor = "#000000"
	}

	attr := map[string]any{
		"editable":      false,
		"pattern":       "null",
		"mode":          "text",
		"icon":          "mdi:fire",
		"friendly_name": h.config.HaFrontforceStatusEntity.FriendlyName,
		"color":         currAvail.Color,
		"text_color":    textColor,
	}
	err := h.updateState(h.config.HaFrontforceStatusEntity.EntityID, attr, currAvail.Description)
	if err != nil {
		return err
	}
	return nil
}

func (h HomeAssistant) UpdateInterventionState(interventions []shared.Intervention) error {
	attr := map[string]any{
		"editable":      false,
		"pattern":       "null",
		"mode":          "text",
		"icon":          "mdi:fire",
		"friendly_name": h.config.HaFrontforceInterventionEntity.FriendlyName,
	}

	var interventionDescriptions []string
	for _, interv := range interventions {
		if interv.Type.Category.Name != "Ziekenwagen" {
			interventionDescriptions = append(interventionDescriptions, interv.Type.Description)
		}
	}
	stateValue := strings.Join(interventionDescriptions, ", ")

	err := h.updateState(h.config.HaFrontforceInterventionEntity.EntityID, attr, stateValue)
	if err != nil {
		return err
	}
	return nil
}

func (h HomeAssistant) updateState(entityId string, attributes map[string]interface{}, value string) error {
	var bearer = "Bearer " + h.config.HaToken

	postBody, err := json.Marshal(map[string]interface{}{
		"state":      value,
		"attributes": attributes,
	})
	if err != nil {
		log.Error().Err(err).Msg("home assistant - failed marshalling state update value")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s/%s", h.config.HaUrl, "api/states", entityId), bytes.NewBuffer(postBody))
	if err != nil {
		log.Error().Err(err).Msg("home assistant - failed creating update state request")
		return err
	}

	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msgf("home assistant - failed updating %s state", entityId)
		return err
	}
	if resp.StatusCode != 200 {
		log.Error().Msgf("home assistant - expected code 200, received: %d", resp.StatusCode)
		return err
	}
	log.Info().Msgf("home assistant - successfully updated %s state", entityId)
	return nil
}

func calcAvailabilityStats(currentAvailability float64) (float64, float64, float64, float64) {
	current := time.Now()
	currentYear := current.Year()
	currentMonth := current.Month()

	start := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, current.Location())
	end := time.Date(currentYear, currentMonth+1, 1, 0, 0, 0, 0, current.Location())
	secondsInMonth := end.Sub(start).Seconds()

	neededSeconds := secondsInMonth * 0.25
	neededHours := neededSeconds / time.Hour.Seconds()

	secondsPast := current.Sub(start).Seconds()
	secondsAvailable := secondsPast * currentAvailability / 100.0
	hoursAvailable := secondsAvailable / time.Hour.Seconds()

	secondsRemaining := end.Add(time.Duration(-1) * time.Second).Sub(current).Seconds()
	maxSecondsAvailable := secondsAvailable + secondsRemaining
	maxPercentage := maxSecondsAvailable / secondsInMonth * 100

	hoursToGo := neededHours - hoursAvailable
	if hoursToGo < 0 {
		hoursToGo = 0
	}

	return hoursAvailable, neededHours, hoursToGo, maxPercentage
}
