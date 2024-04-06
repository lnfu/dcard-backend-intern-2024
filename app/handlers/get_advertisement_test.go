package handlers

import (
	"database/sql"
	"errors"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
	sqlc "github.com/lnfu/dcard-intern/app/models/sqlc"
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

func TestHandler_buildDBParams(t *testing.T) {
	handler := &Handler{}

	tests := []struct {
		name            string
		queryParameters QueryParameters
		expectedParams  sqlc.GetActiveAdvertisementsParams
	}{
		{
			name: "valid (all)",
			queryParameters: QueryParameters{
				Age:      Int32Ptr(25),
				Gender:   StringPtr("M"),
				Country:  StringPtr("TW"),
				Platform: StringPtr("android"),
				Offset:   Int32Ptr(0),
				Limit:    Int32Ptr(5),
			},
			expectedParams: sqlc.GetActiveAdvertisementsParams{
				Age:      sql.NullInt32{Int32: 25, Valid: true},
				Gender:   sql.NullString{String: "M", Valid: true},
				Country:  sql.NullString{String: "TW", Valid: true},
				Platform: sql.NullString{String: "android", Valid: true},
				Offset:   0,
				Limit:    5,
			},
		},
		{
			name: "only age",
			queryParameters: QueryParameters{
				Age: Int32Ptr(30),
			},
			expectedParams: sqlc.GetActiveAdvertisementsParams{
				Age:      sql.NullInt32{Int32: 30, Valid: true},
				Gender:   sql.NullString{Valid: false},
				Country:  sql.NullString{Valid: false},
				Platform: sql.NullString{Valid: false},
				Offset:   0, // 預設值
				Limit:    5, // 預設值
			},
		},
		{
			name: "only gender",
			queryParameters: QueryParameters{
				Gender: StringPtr("M"),
			},
			expectedParams: sqlc.GetActiveAdvertisementsParams{
				Age:      sql.NullInt32{Valid: false},
				Gender:   sql.NullString{String: "M", Valid: true},
				Country:  sql.NullString{Valid: false},
				Platform: sql.NullString{Valid: false},
				Offset:   0, // 預設值
				Limit:    5, // 預設值
			},
		},
		{
			name: "only country",
			queryParameters: QueryParameters{
				Country: StringPtr("TW"),
			},
			expectedParams: sqlc.GetActiveAdvertisementsParams{
				Age:      sql.NullInt32{Valid: false},
				Gender:   sql.NullString{Valid: false},
				Country:  sql.NullString{String: "TW", Valid: true},
				Platform: sql.NullString{Valid: false},
				Offset:   0, // 預設值
				Limit:    5, // 預設值
			},
		},
		{
			name: "only platform",
			queryParameters: QueryParameters{
				Platform: StringPtr("android"),
			},
			expectedParams: sqlc.GetActiveAdvertisementsParams{
				Age:      sql.NullInt32{Valid: false},
				Gender:   sql.NullString{Valid: false},
				Country:  sql.NullString{Valid: false},
				Platform: sql.NullString{String: "android", Valid: true},
				Offset:   0, // 預設值
				Limit:    5, // 預設值
			},
		},
		{
			name: "only offset",
			queryParameters: QueryParameters{
				Offset: Int32Ptr(12),
			},
			expectedParams: sqlc.GetActiveAdvertisementsParams{
				Age:      sql.NullInt32{Valid: false},
				Gender:   sql.NullString{Valid: false},
				Country:  sql.NullString{Valid: false},
				Platform: sql.NullString{Valid: false},
				Offset:   12,
				Limit:    5, // 預設值
			},
		},
		{
			name: "only limit",
			queryParameters: QueryParameters{
				Limit: Int32Ptr(20),
			},
			expectedParams: sqlc.GetActiveAdvertisementsParams{
				Age:      sql.NullInt32{Valid: false},
				Gender:   sql.NullString{Valid: false},
				Country:  sql.NullString{Valid: false},
				Platform: sql.NullString{Valid: false},
				Offset:   0, // 預設值
				Limit:    20,
			},
		},
		{
			name:            "empty",
			queryParameters: QueryParameters{},
			expectedParams: sqlc.GetActiveAdvertisementsParams{
				Age:      sql.NullInt32{Valid: false},
				Gender:   sql.NullString{Valid: false},
				Country:  sql.NullString{Valid: false},
				Platform: sql.NullString{Valid: false},
				Offset:   0, // 預設值
				Limit:    5, // 預設值
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			params := handler.buildDBParams(test.queryParameters)
			if params.Age != test.expectedParams.Age ||
				params.Gender != test.expectedParams.Gender ||
				params.Country != test.expectedParams.Country ||
				params.Platform != test.expectedParams.Platform ||
				params.Offset != test.expectedParams.Offset ||
				params.Limit != test.expectedParams.Limit {
				t.Errorf("expected: %+v, got: %+v", test.expectedParams, params)
			}
		})
	}
}
