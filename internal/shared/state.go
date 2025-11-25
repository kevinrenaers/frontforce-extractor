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
