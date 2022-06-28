package middlewares

import (
	"github.com/piqba/common/pkg/dql"
	"github.com/piqba/common/pkg/httpsrv"
	"net/http"
	"strings"
)

// TODO: ask what is append (owner & space)
// TODO: ask what is requiere space filter

// VerifyTenantHeader  discriminate...
func VerifyTenantHeader(spaceSettings map[string]interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			xWorkSpace := r.Header.Get(WorkSpaceXHeader)
			commonCtx, _ := FromCommonContext(r.Context(), CommonCtxKey)

			switch {
			case spaceSettings[SpaceMultipleProperty] == true:
				if strings.Contains(strings.Join(commonCtx.ClaimCtx.Space.SecureSpacePath, ","), commonCtx.ClaimCtx.Space.Id) {
					// if is multiple we check that the current space is one of the secure space path
					appendSecurePathFilter(commonCtx)
					ctx := NewCommonContext(r.Context(), CommonCtxKey, commonCtx)
					next.ServeHTTP(w, r.WithContext(ctx)) // dispatch the request
					return
				}
			case xWorkSpace == WorkSpaceXHeaderEmpty:
				httpsrv.ResponseError(w, http.StatusBadRequest, "error: missing x-workspace header")
				return
			case xWorkSpace == commonCtx.ClaimCtx.Space.Identifier:
				appendSpaceFilter(commonCtx)
				ctx := NewCommonContext(r.Context(), CommonCtxKey, commonCtx)
				next.ServeHTTP(w, r.WithContext(ctx)) // dispatch the request
				return
			case commonCtx.ClaimCtx.Space.IsRootSpace:
				next.ServeHTTP(w, r.WithContext(r.Context())) // dispatch the request
				return
				// TODO: when is public space do ????
			default:
				httpsrv.ResponseError(w, http.StatusForbidden, "error: you don`t have permissions to view this content")
				return
			}
		})
	}
}

func appendSecurePathFilter(commonCtx *CommonCtx) {
	spaceSecurePathFilter := dql.Filters{
		Type:       "TERMS",
		Field:      "secureSpacePath",
		Value:      commonCtx.ClaimCtx.Space.SecureSpacePath,
		SubFilters: nil,
	}
	rqCtx := NewRequestSearchDefault()
	rqCtx.Filters.Type = "AND"
	rqCtx.Filters.SubFilters = append(rqCtx.Filters.SubFilters, spaceSecurePathFilter)
	commonCtx.RequestSearchCtx = *rqCtx
}

func appendSpaceFilter(commonCtx *CommonCtx) {
	spaceFilter := dql.Filters{
		Type:       "TERM",
		Field:      "space",
		Value:      commonCtx.ClaimCtx.Space.Id,
		SubFilters: nil,
	}
	rqCtx := NewRequestSearchDefault()
	rqCtx.Filters.Type = "AND"
	rqCtx.Filters.SubFilters = append(rqCtx.Filters.SubFilters, spaceFilter)
	commonCtx.RequestSearchCtx = *rqCtx
}
