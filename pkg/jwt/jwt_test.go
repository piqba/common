package jwt

import (
	"errors"
	"testing"
)

var (
	hmacSecret = "secret"
	// For testing propourse, in the future this token could be expired
	tokenToVerify = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI2MjE1MGQ1YTMwZDk5MzFhYmJkNmFkOTQiLCJpZCI6IjYyMTUwZDVhMzBkOTkzMWFiYmQ2YWQ5NCIsImZpcnN0TmFtZSI6IkFkbWluIiwibGFzdE5hbWUiOiJTeXN0ZW0iLCJlbWFpbCI6ImFkbWluQGFkbWluLmNvbSIsInJvbGVzIjpbeyJyb2xlIjoiNjIxNTJlN2Q2ZjgwZDkzYWNlOThhNWY1Iiwic3BhY2UiOiI2MjE1MmU3YTZmODBkOTNhY2U5OGE1ZjMifV0sIm1ldGFkYXRhIjp7fSwibGFuZ3VhZ2UiOiJlcyIsImZ1bGxOYW1lIjoiQWRtaW4gU3lzdGVtIiwib25ib2FyZGluZ0NvbXBsZXRlZCI6ZmFsc2UsInBlcm1pc3Npb25zIjpbIkFETUlOIl0sInNlY3VyZVNwYWNlUGF0aCI6W10sImlhdCI6MTY0NTgyMDIwMSwiZXhwIjoxNjQ1ODI3NDAxfQ.zyQ5XWjm8hpBw2Ud0L1DDMSCU__SM5CS207zrYI5vYY"

	tokenDataTable = []struct {
		name          string
		token         string
		errorExpected bool
		errorMessage  string
	}{
		{
			name:          "valid-token",
			token:         tokenToVerify,
			errorExpected: false,
			errorMessage:  "no invalid token",
		},
		{
			name:          "mal-formed-token",
			token:         "TokenToVerify",
			errorExpected: true,
			errorMessage:  "Mal Formed Token",
		},
		{
			name:          "expired-token",
			token:         tokenToVerify,
			errorExpected: true,
			errorMessage:  "token expired or not valid yet",
		},
	}
)

// Test_TokenValidator test for verification of JWT
func Test_TokenValidator(t *testing.T) {
	for _, tkd := range tokenDataTable {
		if tkd.errorExpected {
			_, err := TokenValidator(tokenToVerify, hmacSecret)
			if errors.Is(err, ErrValidationErrorMalformed) {
				t.Log("Mal Formed Token")
			} else if errors.Is(err, ErrValidationErrorExpiredOrNotValidYet) {
				t.Log("token expired or not valid yet")
			} else if err != nil {
				t.Logf("error: couldn't handle this token %v", err)
			}
		} else {
			t.Log("Token OK")
		}
	}
}
