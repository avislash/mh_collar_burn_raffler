package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/avislash/mh_collar_burn_raffler/metadata"
)

type HoundsMetadataFetcher struct {
	baseURL string
}

func NewHoundsMetadataFetcher(baseURL string) *HoundsMetadataFetcher {
	return &HoundsMetadataFetcher{baseURL}
}

func (hmf *HoundsMetadataFetcher) FetchMetdata(tokenIDs []uint64) ([]metadata.HoundMetadata, error) {
	houndsMetadata := make([]metadata.HoundMetadata, len(tokenIDs))
	for i, tokenID := range tokenIDs {
		var metadata metadata.HoundMetadata
		response, err := http.Get(fmt.Sprintf("%s/%d", hmf.baseURL, tokenID))
		if err != nil {
			return nil, fmt.Errorf("Error fetching metadata for token %d: %w", i, err)
		}

		err = json.NewDecoder(response.Body).Decode(&metadata) //json.Unmarshal(responseData, &metadata)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling metadata: %w", err)
		}
		houndsMetadata[i] = metadata
	}
	return houndsMetadata, nil
}
