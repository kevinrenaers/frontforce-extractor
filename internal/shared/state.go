package shared

import (
	"encoding/json"
	"fmt"

	"renaers.be/frontforce/internal/configuration"
)

type State struct {
	config  configuration.Configuration
	blocks  map[string]json.RawMessage
	Persons []Person
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

	for _, block := range s.blocks {
		pb := personBlock{}
		err := json.Unmarshal(block, &pb)
		if err != nil {
			return fmt.Errorf("frontforce - failed unmarshalling person block: %w", err)
		}
		persons = append(persons, pb.Data.Persons...)
	}

	s.Persons = persons

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
