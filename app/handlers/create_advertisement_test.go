package handlers

import (
	"errors"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

func TestHandler_validateCondition(t *testing.T) {
	handler := Handler{
		genderSet:   mapset.NewSet("M", "F"),
		countrySet:  mapset.NewSet("TW", "US", "JP"),
		platformSet: mapset.NewSet("android", "ios", "web"),
	}

	testCases := []struct {
		name          string
		condition     AdvertisementCondition
		expectedError error
	}{
		{
			name: "valid condition (all)",
			condition: AdvertisementCondition{
				AgeStart: Int32Ptr(20),
				AgeEnd:   Int32Ptr(30),
				Gender:   []string{"M", "F"},
				Country:  []string{"TW", "JP"},
				Platform: []string{"android"},
			},
			expectedError: nil,
		},
		{
			name: "valid condition (partial)",
			condition: AdvertisementCondition{
				AgeStart: nil,
				AgeEnd:   nil,
				Gender:   []string{"M", "F"},
				Country:  []string{"TW", "JP"},
				Platform: []string{"android"},
			},
			expectedError: nil,
		},
		{
			name: "valid condition (partial)",
			condition: AdvertisementCondition{
				AgeStart: Int32Ptr(20),
				AgeEnd:   Int32Ptr(30),
				Gender:   []string{},
				Country:  []string{},
				Platform: []string{},
			},
			expectedError: nil,
		},
		{
			name: "invalid ageStart (zero)",
			condition: AdvertisementCondition{
				AgeStart: Int32Ptr(0),
			},
			expectedError: errors.New("invalid ageStart value (must be 1 ~ 100)"),
		},
		{
			name: "invalid ageStart (negative)",
			condition: AdvertisementCondition{
				AgeStart: Int32Ptr(-1),
			},
			expectedError: errors.New("invalid ageStart value (must be 1 ~ 100)"),
		},
		{
			name: "invalid ageEnd (zero)",
			condition: AdvertisementCondition{
				AgeEnd: Int32Ptr(-1),
			},
			expectedError: errors.New("invalid ageEnd value (must be 1 ~ 100)"),
		},
		{
			name: "invalid ageEnd (negative)",
			condition: AdvertisementCondition{
				AgeEnd: Int32Ptr(-1),
			},
			expectedError: errors.New("invalid ageEnd value (must be 1 ~ 100)"),
		},
		{
			name: "invalid ageEnd (< ageStart)",
			condition: AdvertisementCondition{
				AgeStart: Int32Ptr(50),
				AgeEnd:   Int32Ptr(20),
			},
			expectedError: errors.New("invalid ageEnd value (must be >= ageStart)"),
		},
		{
			name: "invalid gender",
			condition: AdvertisementCondition{
				Gender: []string{"F", "X"},
			},
			expectedError: errors.New("invalid gender value"),
		},
		{
			name: "invalid country",
			condition: AdvertisementCondition{
				Country: []string{"TW", "AA"},
			},
			expectedError: errors.New("invalid country value"),
		},
		{
			name: "invalid platform",
			condition: AdvertisementCondition{
				Platform: []string{"android", "computer"},
			},
			expectedError: errors.New("invalid platform value"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := handler.validateCondition(tc.condition)
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
