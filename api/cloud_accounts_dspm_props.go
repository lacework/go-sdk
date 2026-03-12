package api

import "encoding/json"

// DspmProps represents the props section for DSPM integrations.
//
// The API returns props in two different formats depending on the endpoint:
//   - GET returns props as a JSON object with camelCase keys
//   - PATCH returns props as a JSON-encoded string with UPPER_SNAKE_CASE keys
//
// UnmarshalJSON handles both forms transparently.
type DspmProps struct {
	Dspm *DspmPropsConfig `json:"dspm,omitempty"`
}

func (p *DspmProps) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	// If data is a JSON string, unwrap it first. This handles the PATCH
	// response format where props is a JSON-encoded string.
	if len(data) > 0 && data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		if s == "" {
			return nil
		}
		data = []byte(s)
	}

	// Try camelCase keys first (GET response format).
	// "type plain DspmProps" creates an alias without the custom UnmarshalJSON
	// method, avoiding infinite recursion when calling json.Unmarshal.
	type plain DspmProps
	if err := json.Unmarshal(data, (*plain)(p)); err == nil && p.Dspm != nil {
		return nil
	}

	// Fall back to UPPER_SNAKE_CASE keys (PATCH response format)
	var upper struct {
		Dspm *dspmPropsConfigUpper `json:"DSPM,omitempty"`
	}
	if err := json.Unmarshal(data, &upper); err != nil {
		return err
	}
	if upper.Dspm != nil {
		p.Dspm = upper.Dspm.toConfig()
	}
	return nil
}

// DspmPropsConfig contains DSPM-specific configuration
type DspmPropsConfig struct {
	ScanIntervalHours *int                  `json:"scanIntervalHours,omitempty"`
	MaxDownloadBytes  *int                  `json:"maxDownloadBytes,omitempty"`
	DatastoreFilters  *DspmDatastoreFilters `json:"datastoreFilters,omitempty"`
}

// DspmDatastoreFilters configures which datastores to include/exclude from scanning
type DspmDatastoreFilters struct {
	FilterMode     string   `json:"filterMode"`
	DatastoreNames []string `json:"datastoreNames"`
}

// dspmPropsConfigUpper mirrors DspmPropsConfig with UPPER_SNAKE_CASE JSON tags
// for deserializing PATCH responses.
type dspmPropsConfigUpper struct {
	ScanIntervalHours *int                       `json:"SCAN_INTERVAL_HOURS,omitempty"`
	MaxDownloadBytes  *int                       `json:"MAX_DOWNLOAD_BYTES,omitempty"`
	DatastoreFilters  *dspmDatastoreFiltersUpper `json:"DATASTORE_FILTERS,omitempty"`
}

type dspmDatastoreFiltersUpper struct {
	FilterMode     string   `json:"FILTER_MODE"`
	DatastoreNames []string `json:"DATASTORE_NAMES"`
}

func (u *dspmPropsConfigUpper) toConfig() *DspmPropsConfig {
	cfg := &DspmPropsConfig{
		ScanIntervalHours: u.ScanIntervalHours,
		MaxDownloadBytes:  u.MaxDownloadBytes,
	}
	if u.DatastoreFilters != nil {
		cfg.DatastoreFilters = &DspmDatastoreFilters{
			FilterMode:     u.DatastoreFilters.FilterMode,
			DatastoreNames: u.DatastoreFilters.DatastoreNames,
		}
	}
	return cfg
}

const apiV2DspmStatus = "v2/dspm/status"

// DspmStatusRequest is the request body for the DSPM status endpoint
type DspmStatusRequest struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// UpdateDspmStatus updates the status of a DSPM integration using server token auth.
func (svc *CloudAccountsService) UpdateDspmStatus(serverToken string, status DspmStatusRequest) error {
	return svc.client.RequestEncoderDecoderWithToken(
		"POST", apiV2DspmStatus, serverToken, status, nil,
	)
}
