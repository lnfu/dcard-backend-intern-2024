package utils

import (
	"database/sql"
	"testing"
)

func TestNullInt32FromInt32Pointer(t *testing.T) {
	input := new(int32)
	var got sql.NullInt32
	var expected sql.NullInt32

	// Test case: nil pointer should return NullInt32 with Valid=false
	got = NullInt32FromInt32Pointer(nil)
	expected = sql.NullInt32{Int32: 0, Valid: false}
	if got != expected {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}

	// Test case: valid pointer with non-zero value should return NullInt32 with Valid=true
	*input = 42
	got = NullInt32FromInt32Pointer(input)
	expected = sql.NullInt32{Int32: 42, Valid: true}
	if got != expected {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}

	// Test case: valid pointer with non-zero value should return NullInt32 with Valid=true
	*input = -25
	got = NullInt32FromInt32Pointer(input)
	expected = sql.NullInt32{Int32: -25, Valid: true}
	if got != expected {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}

	// Test case: valid pointer with zero value should return NullInt32 with Valid=true
	*input = 0
	got = NullInt32FromInt32Pointer(input)
	expected = sql.NullInt32{Int32: 0, Valid: true}
	if got != expected {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestNullStringFromStringPointer(t *testing.T) {
	input := new(string)
	var got sql.NullString
	var expected sql.NullString

	// Test case: nil pointer should return NullString with Valid=false
	got = NullStringFromStringPointer(nil)
	expected = sql.NullString{String: "", Valid: false}
	if got != expected {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}

	// Test case: valid pointer with with string (no space) should return NullInt32 with Valid=true
	*input = "rahmat_erwin_asbdullah"
	got = NullStringFromStringPointer(input)
	expected = sql.NullString{String: "rahmat_erwin_asbdullah", Valid: true}
	if got != expected {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}

	// Test case: valid pointer with string (has spaces) should return NullString with Valid=true
	*input = "rahmat erwin asbdullah"
	got = NullStringFromStringPointer(input)
	expected = sql.NullString{String: "rahmat erwin asbdullah", Valid: true}
	if got != expected {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}

	// Test case: valid pointer with string (chinese) should return NullString with Valid=true
	*input = "天哪這根本就是我"
	got = NullStringFromStringPointer(input)
	expected = sql.NullString{String: "天哪這根本就是我", Valid: true}
	if got != expected {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
