package aws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeTrail struct {
	Name                string `json:"Name"`
	TrailARN            string `json:"TrailARN"`
	HomeRegion          string `json:"HomeRegion"`
	IsOrganizationTrail bool   `json:"IsOrganizationTrail"`
	IsMultiRegionTrail  bool   `json:"IsMultiRegionTrail"`
	S3BucketName        string `json:"S3BucketName,omitempty"`
	SnsTopicARN         string `json:"SnsTopicARN,omitempty"`
}

type nopVerboseWriter struct{}

func (nopVerboseWriter) Write(string) {}
func (nopVerboseWriter) Close()       {}

// newTrailPreflight serves ListTrails and GetTrail the way a member account of an AWS Organization
// sees them: GetTrail resolves a trail only by its full ARN. An organization trail appears there as a
// shadow trail, and looking it up by name fails with TrailNotFoundException (LINK-4504).
func newTrailPreflight(t *testing.T, isOrg bool, trails ...fakeTrail) *Preflight {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-amz-json-1.1")
		switch r.Header.Get("X-Amz-Target") {
		case "CloudTrail_20131101.ListTrails":
			list := []map[string]string{}
			for _, trail := range trails {
				list = append(list, map[string]string{
					"Name": trail.Name, "TrailARN": trail.TrailARN, "HomeRegion": trail.HomeRegion,
				})
			}
			require.NoError(t, json.NewEncoder(w).Encode(map[string]any{"Trails": list}))
		case "CloudTrail_20131101.GetTrail":
			var in struct{ Name string }
			require.NoError(t, json.NewDecoder(r.Body).Decode(&in))
			for _, trail := range trails {
				if trail.TrailARN == in.Name {
					require.NoError(t, json.NewEncoder(w).Encode(map[string]any{"Trail": trail}))
					return
				}
			}
			w.Header().Set("X-Amzn-ErrorType", "TrailNotFoundException")
			w.WriteHeader(http.StatusBadRequest)
			require.NoError(t, json.NewEncoder(w).Encode(map[string]string{
				"__type":  "TrailNotFoundException",
				"message": fmt.Sprintf("Unknown trail: %s for the user: 111111111111", in.Name),
			}))
		default:
			t.Errorf("unexpected CloudTrail call %q", r.Header.Get("X-Amz-Target"))
		}
	}))
	t.Cleanup(server.Close)

	return &Preflight{
		awsConfig: aws.Config{
			Region:       "us-east-1",
			BaseEndpoint: aws.String(server.URL),
			Credentials:  aws.AnonymousCredentials{},
		},
		isOrg:         isOrg,
		verboseWriter: nopVerboseWriter{},
	}
}

var (
	controlTowerTrail = fakeTrail{
		Name:                "aws-controltower-BaselineCloudTrail",
		TrailARN:            "arn:aws:cloudtrail:us-west-2:575891105434:trail/aws-controltower-BaselineCloudTrail",
		HomeRegion:          "us-west-2",
		IsOrganizationTrail: true,
		IsMultiRegionTrail:  true,
		S3BucketName:        "aws-controltower-logs-575891105434-us-west-2",
		SnsTopicARN:         "arn:aws:sns:us-west-2:575891105434:aws-controltower-AllConfigNotifications",
	}
	accountTrail = fakeTrail{
		Name:               "account-trail",
		TrailARN:           "arn:aws:cloudtrail:us-east-1:111111111111:trail/account-trail",
		HomeRegion:         "us-east-1",
		IsMultiRegionTrail: true,
		S3BucketName:       "account-trail-bucket",
	}
)

func TestFetchEligibleTrail(t *testing.T) {
	t.Run("a member account skips the organization trail and keeps its own", func(t *testing.T) {
		trail, err := fetchEligibleTrail(newTrailPreflight(t, false, controlTowerTrail, accountTrail))

		require.NoError(t, err)
		require.NotNil(t, trail)
		assert.Equal(t, accountTrail.Name, aws.ToString(trail.Name))
	})

	t.Run("a member account with only the organization trail finds none", func(t *testing.T) {
		trail, err := fetchEligibleTrail(newTrailPreflight(t, false, controlTowerTrail))

		require.NoError(t, err)
		assert.Nil(t, trail, "no eligible trail, so a new one is created")
	})

	t.Run("an organization integration selects the organization trail", func(t *testing.T) {
		trail, err := fetchEligibleTrail(newTrailPreflight(t, true, controlTowerTrail, accountTrail))

		require.NoError(t, err)
		require.NotNil(t, trail)
		assert.Equal(t, controlTowerTrail.Name, aws.ToString(trail.Name))
	})
}

// Control Tower preflight may run from an account that is not the trail's home account, such as a
// delegated administrator, where the trail is also a shadow trail.
func TestFetchControlTowerTrail(t *testing.T) {
	trail, err := fetchControlTowerTrail(newTrailPreflight(t, true, controlTowerTrail))

	require.NoError(t, err)
	require.NotNil(t, trail)
	assert.Equal(t, controlTowerTrail.SnsTopicARN, aws.ToString(trail.SnsTopicARN))
}
