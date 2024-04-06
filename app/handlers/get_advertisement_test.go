package handlers

import (
	"errors"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func TestHandler_validateQueryParameters(t *testing.T) {
	handler := Handler{
		genderSet:   mapset.NewSet("M", "F"),
		countrySet:  mapset.NewSet("TW", "US", "JP"),
		platformSet: mapset.NewSet("android", "ios", "web"),
	}
	testCases := []struct {
		name            string
		queryParameters QueryParameters
		expectedError   error
	}{
		{
			name: "valid query parameters (all)",
			queryParameters: QueryParameters{
				Age:      24,
				Gender:   "M",
				Country:  "TW",
				Platform: "android",
				Offset:   0,
				Limit:    5,
			},
			expectedError: nil,
		},
		{
			name: "valid query parameters (partial)",
			queryParameters: QueryParameters{
				Gender:   "M",
				Platform: "android",
				Limit:    5,
			},
			expectedError: nil,
		},
		{
			name: "invalid age (negative)",
			queryParameters: QueryParameters{
				Age:      -1,
				Gender:   "M",
				Country:  "TW",
				Platform: "android",
				Offset:   0,
				Limit:    5,
			},
			expectedError: errors.New("invalid age value (must be 1 ~ 100)"),
		},
		{
			name: "invalid age (>100)",
			queryParameters: QueryParameters{
				Age:      150,
				Gender:   "M",
				Country:  "TW",
				Platform: "android",
				Offset:   0,
				Limit:    5,
			},
			expectedError: errors.New("invalid age value (must be 1 ~ 100)"),
		},
		{
			name: "invalid gender",
			queryParameters: QueryParameters{
				Age:      24,
				Gender:   "X",
				Country:  "TW",
				Platform: "android",
				Offset:   0,
				Limit:    5,
			},
			expectedError: errors.New("invalid gender value"),
		},
		{
			name: "invalid country",
			queryParameters: QueryParameters{
				Age:      24,
				Gender:   "M",
				Country:  "AA",
				Platform: "android",
				Offset:   0,
				Limit:    5,
			},
			expectedError: errors.New("invalid country value"),
		},
		{
			name: "invalid platform",
			queryParameters: QueryParameters{
				Age:      24,
				Gender:   "M",
				Country:  "TW",
				Platform: "computer",
				Offset:   0,
				Limit:    5,
			},
			expectedError: errors.New("invalid platform value"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := handler.validateQueryParameters(tc.queryParameters)
			if err != nil && tc.expectedError == nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if err == nil && tc.expectedError != nil {
				t.Errorf("expected error: %v, but got nil", tc.expectedError)
				return
			}
			if err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}
