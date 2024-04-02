package main

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	db "github.com/lnfu/dcard-intern/db/sqlc"
)

func TestNonEmptyNullStringFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want sql.NullString
	}{
		{
			"普通字串",
			args{s: "thisisastring"},
			sql.NullString{String: "thisisastring", Valid: true},
		},
		{
			"有空格的字串",
			args{s: "this is a string"},
			sql.NullString{String: "this is a string", Valid: true},
		},
		{
			"空字串",
			args{s: ""},
			sql.NullString{Valid: false},
		},
		{
			"只有空格的字串",
			args{s: "   "},
			sql.NullString{String: "   ", Valid: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NonEmptyNullStringFromString(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NullStringFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullInt32FromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    sql.NullInt32
		wantErr bool
	}{
		{
			"一般數字",
			args{s: "12345"},
			sql.NullInt32{Int32: 12345, Valid: true},
			false,
		},
		{
			"負數",
			args{s: "-54321"},
			sql.NullInt32{Int32: -54321, Valid: true},
			false,
		},
		{
			"零",
			args{s: "0"},
			sql.NullInt32{Int32: 0, Valid: true},
			false,
		},
		{
			"多個零",
			args{s: "00000000000"},
			sql.NullInt32{Int32: 0, Valid: true},
			false,
		},
		{
			"零開頭",
			args{s: "000000000001"},
			sql.NullInt32{Int32: 1, Valid: true},
			false,
		},
		{
			"空字串",
			args{s: ""},
			sql.NullInt32{Valid: false},
			false,
		},
		{
			"只有空格的字串",
			args{s: "     "},
			sql.NullInt32{Valid: false},
			false,
		},
		{
			"非數字",
			args{s: "abc"},
			sql.NullInt32{Valid: false},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NullInt32FromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("NullInt32FromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NullInt32FromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_application_getAdvertisementFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// TODO 完成 test 會用到的環境設定
	dbConnection, _ := sql.Open(dbDriver, dbSource)
	defer dbConnection.Close()
	app := &application{
		errorLogger:     log.New(os.Stderr, color.RedString("ERROR\t"), log.Ldate|log.Ltime|log.Lshortfile),
		infoLoggger:     log.New(os.Stdout, color.BlueString("INFO\t"), log.Ldate|log.Ltime|log.Lshortfile),
		databaseQueries: db.New(dbConnection),
	}

	tests := []struct {
		name         string
		queryParams  map[string]string
		wantAge      sql.NullInt32
		wantGender   sql.NullString
		wantCountry  sql.NullString
		wantPlatform sql.NullString
		wantOffset   int
		wantLimit    int
		wantErr      bool
	}{
		{
			name: "規格範例",
			queryParams: map[string]string{
				"offset":   "10",
				"limit":    "3",
				"age":      "24",
				"gender":   "F",
				"country":  "TW",
				"platform": "ios",
			},
			wantAge:      sql.NullInt32{Int32: 24, Valid: true},
			wantGender:   sql.NullString{String: "F", Valid: true},
			wantCountry:  sql.NullString{String: "TW", Valid: true},
			wantPlatform: sql.NullString{String: "ios", Valid: true},
			wantOffset:   10,
			wantLimit:    3,
			wantErr:      false,
		},
		{
			name: "預設 limit 和 offset",
			queryParams: map[string]string{
				"age":      "24",
				"gender":   "F",
				"country":  "TW",
				"platform": "ios",
			},
			wantAge:      sql.NullInt32{Int32: 24, Valid: true},
			wantGender:   sql.NullString{String: "F", Valid: true},
			wantCountry:  sql.NullString{String: "TW", Valid: true},
			wantPlatform: sql.NullString{String: "ios", Valid: true},
			wantOffset:   0,
			wantLimit:    5,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			query := req.URL.Query()
			for key, value := range tt.queryParams {
				query.Set(key, value)
			}
			req.URL.RawQuery = query.Encode()

			responseRecorder := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(responseRecorder)
			ctx.Request = req

			age, gender, country, platform, offset, limit, err := app.getAdvertisementFilters(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if age != tt.wantAge {
				t.Errorf("got age %v, want %v", age, tt.wantAge)
			}
			if gender != tt.wantGender {
				t.Errorf("got gender %v, want %v", gender, tt.wantGender)
			}
			if country != tt.wantCountry {
				t.Errorf("got country %v, want %v", country, tt.wantCountry)
			}
			if platform != tt.wantPlatform {
				t.Errorf("got platform %v, want %v", platform, tt.wantPlatform)
			}
			if offset != tt.wantOffset {
				t.Errorf("got offset %d, want %d", offset, tt.wantOffset)
			}
			if limit != tt.wantLimit {
				t.Errorf("got limit %d, want %d", limit, tt.wantLimit)
			}
		})
	}
}
