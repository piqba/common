package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
)

type MapClaims jwt.MapClaims

// JWT is the struct for parse json from token ms-auth (new & old)
type JWT struct {
	BsonId              string        `json:"_id"`
	Id                  string        `json:"id"`
	FirstName           string        `json:"firstName"`
	LastName            string        `json:"lastName"`
	Email               string        `json:"email"`
	Phone               string        `json:"phone"`
	Roles               []interface{} `json:"roles"`
	Metadata            Metadata      `json:"metadata"`
	Language            string        `json:"language"`
	FullName            string        `json:"fullName"`
	IdNumber            string        `json:"idNumber"`
	OnboardingCompleted bool          `json:"onboardingCompleted"`
	Permissions         []string      `json:"permissions"`
	Space               Space         `json:"space"`
	Iat                 int           `json:"iat"`
}
type Metadata struct {
}
type Space struct {
	Id              string   `json:"_id"`
	Identifier      string   `json:"identifier"`
	IsPublicSpace   bool     `json:"isPublicSpace"`
	IsRootSpace     bool     `json:"isRootSpace"`
	SecureSpacePath []string `json:"secureSpacePath"`
}

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
