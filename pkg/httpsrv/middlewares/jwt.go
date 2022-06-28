package middlewares

import (
	"errors"
	"github.com/piqba/common/pkg/httpsrv"
	"github.com/piqba/common/pkg/jwt"
	"net/http"
	"strings"
)

// VerifyJWT ...
func VerifyJWT(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerAuthorization := r.Header.Get(AuthorizationHeader)

			if !strings.HasPrefix(headerAuthorization, "Bearer ") {
				httpsrv.ResponseError(w, http.StatusBadRequest, jwt.ErrBearerTokenFormat.Error())
				return
			}

			tokenToVerify := strings.Split(headerAuthorization, " ")[1]
			claimsJWT, err := jwt.TokenValidator(tokenToVerify, secret)
			switch {
			case errors.Is(err, jwt.ErrValidationErrorMalformed):
				httpsrv.ResponseError(w, http.StatusForbidden, jwt.ErrValidationErrorMalformed.Error())
				return
			case errors.Is(err, jwt.ErrValidationErrorExpiredOrNotValidYet):
				httpsrv.ResponseError(w, http.StatusForbidden, jwt.ErrValidationErrorExpiredOrNotValidYet.Error())
				return
			case err != nil:
				httpsrv.ResponseError(w, http.StatusForbidden, "couldn't handle this token")
				return
			}

			claimCtx := extractClaims(claimsJWT)
			ctx := NewCommonContext(r.Context(), CommonCtxKey, &CommonCtx{
				ClaimCtx: *claimCtx,
			})
			next.ServeHTTP(w, r.WithContext(ctx)) // dispatch the request
		})
	}
}

func extractClaims(claims *jwt.JWT) *ClaimCtx {
	return &ClaimCtx{
		permissions: claims.Permissions,
		email:       claims.Email,
		userId:      claims.Id,
		Space: Space{
			Id:              claims.Space.Id,
			Identifier:      claims.Space.Identifier,
			IsPublicSpace:   claims.Space.IsPublicSpace,
			IsRootSpace:     claims.Space.IsRootSpace,
			SecureSpacePath: claims.Space.SecureSpacePath,
		},
	}
}
