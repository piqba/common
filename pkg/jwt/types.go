package jwt

import "errors"

var (
	// For HMAC signing method, the key can be any []byte. It is recommended to generate
	// a key using crypto/rand or something equivalent. You need the same key for signing
	// and validating.
	//hmacSecret                             = []byte(os.Getenv("APP_MS_AUTH_JWT_SECRET"))

	// ErrBearerTokenFormat when Format is Authorization: Bearer [token]
	ErrBearerTokenFormat = errors.New("error: Format is Authorization: Bearer [token]")
	// ErrValidationErrorMalformed when That's not even a token
	ErrValidationErrorMalformed = errors.New("error: That's not even a token")
	// ErrValidationErrorExpiredOrNotValidYet when  Token is either expired or not active yet
	ErrValidationErrorExpiredOrNotValidYet = errors.New("error: Token is either expired or not active yet")
)
