package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// SearchCtx ...
func SearchCtx() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			buf, _ := ioutil.ReadAll(r.Body)
			bodyCopy := ioutil.NopCloser(bytes.NewBuffer(buf))
			bodyReaderCloser := ioutil.NopCloser(bytes.NewBuffer(buf))

			r.Body = bodyReaderCloser
			ctx := NewCommonContext(
				r.Context(),
				CommonCtxKey,
				injectReqSearchOnCtx(
					bodyCopy,
					r,
				),
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func injectReqSearchOnCtx(bodyCopy io.ReadCloser, req *http.Request) *CommonCtx {
	parsedReq := &RequestSearchCtx{}
	common, _ := FromCommonContext(req.Context(), CommonCtxKey)
	json.NewDecoder(bodyCopy).Decode(&parsedReq)

	common.RequestSearchCtx.Page = parsedReq.DefaultPage(parsedReq.Page, 1)
	common.RequestSearchCtx.Size = parsedReq.DefaultSize(parsedReq.Size, 20)
	common.RequestSearchCtx.Sort = parsedReq.DefaultSort(parsedReq.Sort, "desc")

	common.RequestSearchCtx.Search = parsedReq.Search

	common.RequestSearchCtx.Filters.SubFilters = append(
		common.RequestSearchCtx.Filters.SubFilters,
		parsedReq.Filters.SubFilters...,
	)
	return common
}
