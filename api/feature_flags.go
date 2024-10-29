package api

import (
	"fmt"
)

type FeatureFlagsService struct {
	client *Client
}

type FeatureFlag string

type FeatureFlags struct {
	Flags []FeatureFlag `json:"flags,omitempty"`
}

type FeatureFlagsResponse struct {
	Data FeatureFlags `json:"data"`
}

func (svc *FeatureFlagsService) GetFeatureFlagsMatchingPrefix(prefix string) (
	response FeatureFlagsResponse, err error,
) {
	apiPath := fmt.Sprintf("%s/%s", apiV2FeatureFlags, prefix)
	err = svc.client.RequestDecoder("GET", apiPath, nil, &response)
	return
}
