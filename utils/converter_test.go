package utils

import (
	"database/sql"
	"reflect"
	"testing"
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
