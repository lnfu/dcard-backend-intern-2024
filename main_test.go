package main

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestNullStringFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want sql.NullString
	}{
		{
			"test normal string",
			args{s: "thisisastring"},
			sql.NullString{String: "thisisastring", Valid: true},
		},
		{
			"test string containing spaces",
			args{s: "this is a string"},
			sql.NullString{String: "this is a string", Valid: true},
		},
		{
			"test empty string",
			args{s: ""},
			sql.NullString{Valid: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NullStringFromString(tt.args.s); !reflect.DeepEqual(got, tt.want) {
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
			"test positive",
			args{s: "12345"},
			sql.NullInt32{Int32: 12345, Valid: true},
			false,
		},
		{
			"test negative",
			args{s: "-54321"},
			sql.NullInt32{Int32: -54321, Valid: true},
			false,
		},
		{
			"test zero",
			args{s: "0"},
			sql.NullInt32{Int32: 0, Valid: true},
			false,
		},
		{
			"test multiple zeros",
			args{s: "00000000000"},
			sql.NullInt32{Int32: 0, Valid: true},
			false,
		},
		{
			"test non number",
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
