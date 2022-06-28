package jwt

// TODO: Improve this test
import (
	"errors"
	"testing"
)

var (
	hmacSecret = "secret"
	// For testing propourse, in the future this token could be expired
	tokenToVerify = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQi10sImlhdCI6MTY0NzUyNzcxNywiZXhwIjoxNjQ3NjE0MTE3fQ.o9XL3dlEiZM3A9LM8sNPPj5TXXQ_iGBC7UCvwBb9sNI"

	newToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJYzc2ZGZiYTIzZWY0MWJiZDBlIl0sImlhdCI6MTY0NjI1MTg5MH0.pkUpA1mCRx0KHiuheQYHrJjqGxdGb_IMxpmeys6RCTk"

	tokenDataTable = []struct {
		name          string
		token         string
		errorExpected bool
		errorMessage  string
	}{
		{
			name:          "valid-token",
			token:         newToken,
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
