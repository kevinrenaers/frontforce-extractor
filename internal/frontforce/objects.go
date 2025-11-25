package frontforce

import (
	"encoding/json"
	"slices"

	"renaers.be/frontforce/internal/shared"
)

type frontforceMessage []map[string]json.RawMessage

func (f frontforceMessage) FetchRelevantBlocks() map[string]json.RawMessage {
	result := make(map[string]json.RawMessage)

	for _, blocks := range f {
		for blockID, block := range blocks {
			if !slices.Contains(shared.RelevantBlockIDs, blockID) {
				continue
			}
			result[blockID] = block
		}
	}

	return result
}
