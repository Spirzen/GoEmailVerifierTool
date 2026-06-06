package main

import (
	"testing"
)

func TestEmailRegex(t *testing.T) {
	valid := []string{
		"user@example.com",
		"user.name+tag@sub.example.co.uk",
	}
	invalid := []string{
		"not-an-email",
		"@example.com",
		"user@",
		"user@.com",
	}

	for _, email := range valid {
		if !emailRegex.MatchString(email) {
			t.Errorf("expected valid: %q", email)
		}
	}
	for _, email := range invalid {
		if emailRegex.MatchString(email) {
			t.Errorf("expected invalid: %q", email)
		}
	}
}

func TestParseDomain(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"user@gmail.com", "gmail.com", false},
		{"gmail.com", "gmail.com", false},
		{"user@", "", true},
	}

	for _, tt := range tests {
		got, err := ParseDomain(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseDomain(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseDomain(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
