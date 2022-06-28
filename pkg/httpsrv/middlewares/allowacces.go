package middlewares

import (
	"github.com/piqba/common/pkg/httpsrv"
	"net/http"
)

type AclWords string

const (
	FullAcl  AclWords = "*"
	GuestAcl AclWords = "guest"
	AdminAcl AclWords = "admin"
)

// AllowAccess ...
func AllowAccess(permissionsToAllow ...AclWords) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			common, _ := FromCommonContext(r.Context(), CommonCtxKey)
			for _, pAccess := range permissionsToAllow {
				if pAccess == FullAcl {
					next.ServeHTTP(w, r)
					return
				}
				for _, permissionFromClaim := range common.ClaimCtx.permissions {
					if permissionFromClaim == string(pAccess) {
						next.ServeHTTP(w, r)
						return
					}
					httpsrv.ResponseError(w, http.StatusForbidden, "error: you don`t have permissions to view this content")
					return
				}
				httpsrv.ResponseError(w, http.StatusForbidden, "error: you don`t have permissions to view this content")
				return
			}
			next.ServeHTTP(w, r) // dispatch the request
		})
	}
}
