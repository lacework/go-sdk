package api_test

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/lacework/go-sdk/api"
	"github.com/stretchr/testify/assert"
)

// TestWindowedSearch test basic functionality of windowedSearch
// windowed search takes a function of type search and filter that implements
// searchable filter
func TestWindowedSearch(t *testing.T) {
	var (
		now          = time.Now().UTC()
		before       = now.AddDate(0, 0, -7) // last 7 days
		testResponse mockSearchResponse
		filter       = mockSearchFilter{
			SearchFilter: api.SearchFilter{
				Filters: []api.Filter{{
					Expression: "eq",
					Field:      "urn",
					Value:      "text",
				}},
				TimeFilter: &api.TimeFilter{
					StartTime: &before,
					EndTime:   &now,
				},
			},
			Value: "MOCK",
		}
	)
	searchCounter = 0
	err := api.WindowedSearch(mockSearch, api.V2ApiMaxSearchWindowDays, api.V2ApiMaxSearchHistoryDays, &testResponse, &filter)
	assert.NoError(t, err)
	assert.Equal(t, searchCounter, 3)
	assert.NotEmpty(t, testResponse.Data)
	assert.Equal(t, testResponse.Data[0].ID, "MOCK_1")
	assert.Equal(t, testResponse.Data[0].Value, "EXAMPLE_VALUE_1")

	// startTime from current Date should be less than max history
	timeDifference := int(math.RoundToEven(time.Now().Sub(*filter.GetTimeFilter().StartTime).Hours() / 24))
	assert.Less(t, timeDifference, api.V2ApiMaxSearchHistoryDays)
}

// TestWindowedSearchMaxHistory test max history is not exceeded
func TestWindowedSearchMaxHistory(t *testing.T) {
	var (
		// set max history and window size. searchCounter should not exceed 1
		maxWindow    = 5
		maxHistory   = 10
		now          = time.Now().UTC()
		before       = now.AddDate(0, 0, -5) // last 5 days
		testResponse mockSearchResponse
		filter       = mockSearchFilter{
			SearchFilter: api.SearchFilter{
				TimeFilter: &api.TimeFilter{
					StartTime: &before,
					EndTime:   &now,
				},
			},
			Value: "MOCK",
		}
	)

	searchCounter = 0
	err := api.WindowedSearch(mockSearch, maxWindow, maxHistory, &testResponse, &filter)
	assert.NoError(t, err)

	// search counter should equal 1. No data is found in search
	assert.Equal(t, searchCounter, 1)
	assert.Empty(t, testResponse.Data)

	// startTime from current date should be less or equal to max history
	timeDifference := int(math.RoundToEven(time.Now().Sub(*filter.GetTimeFilter().StartTime).Hours() / 24))
	assert.Equal(t, timeDifference, maxHistory)
}

// TestWindowedSearchMaxHistory test max window is greater than history
// returns error
func TestWindowedSearchHistory(t *testing.T) {
	var (
		// set max history and window size. Window cannot exceed history size
		maxWindow    = 10
		maxHistory   = 5
		now          = time.Now().UTC()
		before       = now.AddDate(0, 0, -10) // last 10 days
		testResponse mockSearchResponse
		filter       = mockSearchFilter{
			SearchFilter: api.SearchFilter{
				TimeFilter: &api.TimeFilter{
					StartTime: &before,
					EndTime:   &now,
				},
			},
			Value: "MOCK",
		}
	)

	err := api.WindowedSearch(mockSearch, maxWindow, maxHistory, &testResponse, &filter)
	assert.ErrorContains(t, err, "window size cannot be greater than max history")
}

// TestWindowedSearchMaxHistory test history is not divisible by window that final search only searches remainder
// the final search only searches remainder
func TestWindowedSearchHistoryRemainder(t *testing.T) {
	var (
		// set max history and window size. searchCounter should not exceed 1
		maxWindow    = 3
		maxHistory   = 5
		now          = time.Now().UTC()
		before       = now.AddDate(0, 0, -3) // last 3 days
		testResponse mockSearchResponse
		filter       = mockSearchFilter{
			SearchFilter: api.SearchFilter{
				TimeFilter: &api.TimeFilter{
					StartTime: &before,
					EndTime:   &now,
				},
			},
			Value: "MOCK",
		}
	)

	searchCounter = 0
	err := api.WindowedSearch(mockSearch, maxWindow, maxHistory, &testResponse, &filter)
	assert.NoError(t, err)

	// search counter should equal 1. No data is found in search
	assert.Equal(t, searchCounter, 1)
	assert.Empty(t, testResponse.Data)

	// the time difference on final search should adjust as to not exceed the max history
	finalSearchEndTime := filter.GetTimeFilter().EndTime
	finalSearchStartTime := filter.GetTimeFilter().StartTime
	timeDifference := int(math.RoundToEven(finalSearchEndTime.Sub(*finalSearchStartTime).Hours() / 24))
	assert.Equal(t, 2, timeDifference)

	// startTime from current date should be less than max history
	maxTimeDifference := int(math.RoundToEven(filter.GetTimeFilter().StartTime.Sub(time.Now()).Hours() / 24))
	assert.Less(t, maxTimeDifference, maxHistory)
}

func mockSearch(response interface{}, filters api.SearchableFilter) error {
	// simulate the search item not being found in first windows
	mockResponse := mockSearchResponseJson
	if searchCounter < 3 {
		mockResponse = mockSearchResponseEmptyJson
		searchCounter++
	}

	err := json.Unmarshal([]byte(mockResponse), &response)
	if err != nil {
		return err
	}
	return nil
}

var searchCounter int

type mockSearchResponse struct {
	Data []mockSearchResponseData `json:"data"`
}

type mockSearchResponseData struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

func (m mockSearchResponse) GetDataLength() int {
	return len(m.Data)
}

type mockSearchFilter struct {
	api.SearchFilter
	Value string `json:"val"`
}

func (m mockSearchFilter) GetTimeFilter() *api.TimeFilter {
	return m.TimeFilter

}

func (m mockSearchFilter) SetStartTime(t *time.Time) {
	m.TimeFilter.StartTime = t
}

func (m mockSearchFilter) SetEndTime(t *time.Time) {
	m.TimeFilter.EndTime = t
}

var mockSearchResponseJson = `
{
  "data": [
    {
      "id": "MOCK_1",
      "value": "EXAMPLE_VALUE_1" 
    },
    {
      "id": "MOCK_2",
      "value": "EXAMPLE_VALUE_2" 
    }
  ]
}
`

var mockSearchResponseEmptyJson = `
{
  "data": []
}
`
