package handlers

import (
	"errors"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func Int32Ptr(i int32) *int32    { return &i }
func StringPtr(s string) *string { return &s }

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
				Age:      Int32Ptr(25),
				Gender:   StringPtr("M"),
				Country:  StringPtr("TW"),
				Platform: StringPtr("android"),
				Offset:   Int32Ptr(0),
				Limit:    Int32Ptr(5),
			},
			expectedError: nil,
		},
		{
			name: "valid query parameters (partial)",
			queryParameters: QueryParameters{
				Age:      nil,
				Gender:   StringPtr("M"),
				Country:  nil,
				Platform: StringPtr("android"),
				Offset:   nil,
				Limit:    Int32Ptr(5),
			},
			expectedError: nil,
		},
		{
			name: "valid query parameters (partial)",
			queryParameters: QueryParameters{
				Age:      Int32Ptr(25),
				Gender:   nil,
				Country:  StringPtr("TW"),
				Platform: nil,
				Offset:   Int32Ptr(0),
				Limit:    nil,
			},
			expectedError: nil,
		},
		{
			name: "valid query parameters (empty)",
			queryParameters: QueryParameters{
				Age:      nil,
				Gender:   nil,
				Country:  nil,
				Platform: nil,
				Offset:   nil,
				Limit:    nil,
			},
			expectedError: nil,
		},
		{
			name: "invalid age (negative)",
			queryParameters: QueryParameters{
				Age: Int32Ptr(-1),
			},
			expectedError: errors.New("invalid age value (must be 1 ~ 100)"),
		},
		{
			name: "invalid age (zero)",
			queryParameters: QueryParameters{
				Age: Int32Ptr(0),
			},
			expectedError: errors.New("invalid age value (must be 1 ~ 100)"),
		},
		{
			name: "invalid age (> 100)",
			queryParameters: QueryParameters{
				Age: Int32Ptr(101),
			},
			expectedError: errors.New("invalid age value (must be 1 ~ 100)"),
		},
		{
			name: "invalid gender",
			queryParameters: QueryParameters{
				Gender: StringPtr("X"),
			},
			expectedError: errors.New("invalid gender value"),
		},
		{
			name: "invalid country",
			queryParameters: QueryParameters{
				Country: StringPtr("AA"),
			},
			expectedError: errors.New("invalid country value"),
		},
		{
			name: "invalid platform",
			queryParameters: QueryParameters{
				Platform: StringPtr("computer"),
			},
			expectedError: errors.New("invalid platform value"),
		},
		{
			name: "invalid offset (negative)",
			queryParameters: QueryParameters{
				Offset: Int32Ptr(-1),
				Limit:  Int32Ptr(5),
			},
			expectedError: errors.New("invalid offset value (must be >= 0)"),
		},
		{
			name: "invalid limit (zero)",
			queryParameters: QueryParameters{
				Offset: Int32Ptr(0),
				Limit:  Int32Ptr(0),
			},
			expectedError: errors.New("invalid limit value (must be 1 ~ 100)"),
		},
		{
			name: "invalid limit (negative)",
			queryParameters: QueryParameters{
				Offset: Int32Ptr(0),
				Limit:  Int32Ptr(-1),
			},
			expectedError: errors.New("invalid limit value (must be 1 ~ 100)"),
		},
		{
			name: "invalid limit (> 100)",
			queryParameters: QueryParameters{
				Offset: Int32Ptr(0),
				Limit:  Int32Ptr(101),
			},
			expectedError: errors.New("invalid limit value (must be 1 ~ 100)"),
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
