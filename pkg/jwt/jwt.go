package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
)

// TokenValidator it`s used to validate jwt-token and return map[string]interface (jwt.MapClaims)
func TokenValidator(tokenFromHeader, hmacSecret string) (jwt.MapClaims, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenFromHeader, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error: Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSecret, nil
	})
	var claims jwt.MapClaims
	if mapClaims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claims = mapClaims
	}
	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return nil, ErrValidationErrorMalformed
		}
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return nil, ErrValidationErrorExpiredOrNotValidYet
		}
		return nil, errors.New(fmt.Sprintf("error: couldn't handle this token %v", err))
	}

	return claims, nil
}
