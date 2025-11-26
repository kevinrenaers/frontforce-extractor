package frontforce

import (
	"encoding/json"
	"slices"
)

type frontforceMessage []map[string]json.RawMessage

func (f frontforceMessage) FetchRelevantBlocks(relevantBlocks []string) map[string]json.RawMessage {
	result := make(map[string]json.RawMessage)

	for _, blocks := range f {
		for blockID, block := range blocks {
			if !slices.Contains(relevantBlocks, blockID) {
				continue
			}
			result[blockID] = block
		}
	}

	return result
}
