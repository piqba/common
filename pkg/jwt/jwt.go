package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
)

// TokenValidator it`s used to validate jwt-token and return map[string]interface (jwt.MapClaims)
func TokenValidator(tokenFromHeader, hmacSecret string) (*JWT, error) {
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
		return []byte(hmacSecret), nil
	})
	//var claims jwt.MapClaims
	var jwtObj JWT
	if mapClaims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// TODO: jsonMarshal mapClaims
		result, err := fromMapClaimsToJwtObj(mapClaims)
		if err != nil {
			return nil, err
		}
		jwtObj = *result
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

	return &jwtObj, nil
}

// fromMapClaimsToJwtObj ...
func fromMapClaimsToJwtObj(mapClaims jwt.MapClaims) (*JWT, error) {
	var jwtObj JWT
	bytes, err := json.Marshal(&mapClaims)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &jwtObj)
	if err != nil {
		return nil, err
	}
	return &jwtObj, nil
}
