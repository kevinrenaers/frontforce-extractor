package shared

import (
	"encoding/json"
	"fmt"
	"slices"

	"renaers.be/frontforce/internal/configuration"
)

const totalAmountPersons = 49.0

var (
	personBlockIDs       = []string{"1411", "1414", "1415", "1416", "1417", "1418", "1419"}
	interventionsBlockID = "1421"
	vehiclesBlockID      = "1422"
	RelevantBlockIDs     = append(append(personBlockIDs, interventionsBlockID), vehiclesBlockID)
	possibleStatuses     = []string{
		"Niet beschikbaar",
		"Niet beschikbaar",
		"Beschikbaar Werk",
		"Beschikbaar snel",
		"Beschikbaar traag",
		"Dienstopdracht",
		"Reserve",
		"Kazerne Operationeel",
		"Kazerne Administratief",
		"Officier van dienst",
		"Kazerne beschikbaar",
		"Ambulance",
		"Officier beschikbaar",
		"Niet Operationeel",
		"Interventie",
		"Elders beschikbaar",
	}
)

type State struct {
	config        configuration.Configuration
	blocks        map[string]json.RawMessage
	Persons       []Person
	Interventions []Intervention
	Vehicles      []Vehicle
}

func NewState(config configuration.Configuration) *State {
	return &State{
		config: config,
		blocks: make(map[string]json.RawMessage),
	}
}

func (s *State) Update(relevantBlocks map[string]json.RawMessage) error {
	var persons []Person

	for blockID, block := range relevantBlocks {
		s.blocks[blockID] = block
	}

	for blockID, block := range s.blocks {
		if !slices.Contains(personBlockIDs, blockID) {
			continue
		}

		pb := personBlock{}
		err := json.Unmarshal(block, &pb)
		if err != nil {
			return fmt.Errorf("frontforce - failed unmarshalling person block: %w", err)
		}
		persons = append(persons, pb.Data.Persons...)
	}

	s.Persons = persons

	interventionsBlock, ok := s.blocks[interventionsBlockID]
	if ok {
		ib := interventionBlock{}
		err := json.Unmarshal(interventionsBlock, &ib)
		if err != nil {
			return fmt.Errorf("frontforce - failed unmarshalling intervention block: %w", err)
		}
		s.Interventions = ib.Data
	}

	vehiclesBlock, ok := s.blocks[vehiclesBlockID]
	if ok {
		vb := vehicleBlock{}
		err := json.Unmarshal(vehiclesBlock, &vb)
		if err != nil {
			return fmt.Errorf("frontforce - failed unmarshalling vehicles block: %w", err)
		}
		var relevantVehicles []Vehicle
		for _, vehicle := range vb.Data {
			if slices.Contains(s.config.FrontforceVehicleCodes, vehicle.Name) {
				relevantVehicles = append(relevantVehicles, vehicle)
			}
		}
		s.Vehicles = relevantVehicles
	}

	return nil
}

func (s State) GetPerson(id int) Person {
	for _, p := range s.Persons {
		if p.ID == id {
			return p
		}
	}

	return Person{
		UnavailabilityCode: UnavailabilityCode{
			ID:          -1,
			Code:        "NB",
			Description: "Niet beschikbaar",
			Color:       "#dbdbdb",
			IsAvailable: false,
		},
	}
}

func (s State) CalcAvailabilityStats() map[string]any {
	stats := make(map[string]any)

	var amountCallable float64

	statusStats := make(map[string]int)
	for _, status := range possibleStatuses {
		statusStats[status] = 0
	}

	for _, p := range s.Persons {
		code := p.UnavailabilityCode.Description
		statusStats[code]++

		if p.UnavailabilityCode.IsAvailable {
			amountCallable++
		}
	}

	available := float64(len(s.Persons))
	unavailable := totalAmountPersons - available
	stats["post_availability_percentage"] = available / totalAmountPersons * 100
	stats["post_availability_count"] = available
	stats["post_callable_count"] = amountCallable
	stats["post_unavailability_count"] = unavailable
	for status, count := range statusStats {
		stats[status] = count
	}

	return stats
}
