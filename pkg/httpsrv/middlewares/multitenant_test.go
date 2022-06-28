package middlewares

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gitlab.com/dfl-go-pkg/common/pkg/httpsrv"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	indentifierOK  = "identifierValue"
	indentifierBad = "bad"
	spaceSettings  = httpsrv.M{}
	claimTest      = ClaimCtx{
		permissions: []string{"ADMIN"},
		email:       "admin@mail.com",
		Space: Space{
			Id:              "idSpace",
			Identifier:      "identifierValue",
			IsPublicSpace:   true,
			IsRootSpace:     false,
			SecureSpacePath: []string{"61e881c76dfba23ef41bbd0e", "idSpace"},
		},
		userId: "61e881d06dfba23ef41bbd15",
	}
)

func TestMultiTenant_BadIdentifier_NoRoot(t *testing.T) {
	assertions := assert.New(t)
	spaceSettings["multiple"] = false
	want := httpsrv.M{
		"message": "error: you don`t have permissions to view this content",
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	req.Header.Set(WorkSpaceXHeader, indentifierBad)
	ctx := NewCommonContext(req.Context(), CommonCtxKey, &CommonCtx{
		RequestSearchCtx: RequestSearchCtx{},
		ClaimCtx:         claimTest,
	})
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	h := Chain(
		nextHandler,
		VerifyTenantHeader(spaceSettings),
	)
	h.ServeHTTP(w, req.WithContext(ctx))

	if w.Code != http.StatusForbidden {
		t.Errorf("Status code spected %d, getter %d", http.StatusForbidden, w.Code)
	}
	got := httpsrv.M{}
	err := json.NewDecoder(w.Body).Decode(&got)
	if err != nil {
		assertions.Error(err)
	}
	assert.Equal(t, want, got)
}
func TestMultiTenant_EmptyWorkSpaceXHeader(t *testing.T) {
	assertions := assert.New(t)
	spaceSettings["multiple"] = false
	want := httpsrv.M{
		"message": "error: missing x-workspace header",
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	req.Header.Set(WorkSpaceXHeader, WorkSpaceXHeaderEmpty)
	ctx := NewCommonContext(req.Context(), CommonCtxKey, &CommonCtx{
		RequestSearchCtx: RequestSearchCtx{},
		ClaimCtx:         claimTest,
	})
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	h := Chain(
		nextHandler,
		VerifyTenantHeader(spaceSettings),
	)
	h.ServeHTTP(w, req.WithContext(ctx))

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code spected %d, getter %d", http.StatusBadRequest, w.Code)
	}
	got := httpsrv.M{}
	err := json.NewDecoder(w.Body).Decode(&got)
	if err != nil {
		assertions.Error(err)
	}
	assert.Equal(t, want, got)
}
func TestMultiTenant_OK(t *testing.T) {
	assertions := assert.New(t)
	spaceSettings["multiple"] = false
	want := httpsrv.M{}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	req.Header.Set(WorkSpaceXHeader, indentifierOK)
	ctx := NewCommonContext(req.Context(), CommonCtxKey, &CommonCtx{
		RequestSearchCtx: RequestSearchCtx{},
		ClaimCtx:         claimTest,
	})
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	h := Chain(
		nextHandler,
		VerifyTenantHeader(spaceSettings),
	)
	h.ServeHTTP(w, req.WithContext(ctx))

	if w.Code != http.StatusOK {
		t.Errorf("Status code spected %d, getter %d", http.StatusOK, w.Code)
	}
	got := httpsrv.M{}
	err := json.NewDecoder(w.Body).Decode(&got)
	if err != nil {
		assertions.Error(err)
	}
	assert.Equal(t, want, got)
}
func TestMultiTenant_MutilSpace_AppendSecurePath(t *testing.T) {
	assertions := assert.New(t)
	spaceSettings["multiple"] = true
	want := httpsrv.M{}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	req.Header.Set(WorkSpaceXHeader, indentifierOK)
	ctx := NewCommonContext(req.Context(), CommonCtxKey, &CommonCtx{
		RequestSearchCtx: RequestSearchCtx{},
		ClaimCtx:         claimTest,
	})
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		common, _ := FromCommonContext(r.Context(), CommonCtxKey)
		assert.Equal(t, common.RequestSearchCtx.Filters.SubFilters[0].Field, "secureSpacePath")
	})

	h := Chain(
		nextHandler,
		VerifyTenantHeader(spaceSettings),
	)
	h.ServeHTTP(w, req.WithContext(ctx))

	if w.Code != http.StatusOK {
		t.Errorf("Status code spected %d, getter %d", http.StatusOK, w.Code)
	}
	got := httpsrv.M{}
	err := json.NewDecoder(w.Body).Decode(&got)
	if err != nil {
		assertions.Error(err)
	}
	assert.Equal(t, want, got)
}
func TestMultiTenant_MutilSpace_AppendSpace(t *testing.T) {
	assertions := assert.New(t)
	spaceSettings["multiple"] = false
	want := httpsrv.M{}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	req.Header.Set(WorkSpaceXHeader, indentifierOK)
	ctx := NewCommonContext(req.Context(), CommonCtxKey, &CommonCtx{
		RequestSearchCtx: RequestSearchCtx{},
		ClaimCtx:         claimTest,
	})
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		common, _ := FromCommonContext(r.Context(), CommonCtxKey)
		assert.Equal(t, common.RequestSearchCtx.Filters.SubFilters[0].Field, "space")
	})

	h := Chain(
		nextHandler,
		VerifyTenantHeader(spaceSettings),
	)
	h.ServeHTTP(w, req.WithContext(ctx))

	if w.Code != http.StatusOK {
		t.Errorf("Status code spected %d, getter %d", http.StatusOK, w.Code)
	}
	got := httpsrv.M{}
	err := json.NewDecoder(w.Body).Decode(&got)
	if err != nil {
		assertions.Error(err)
	}
	assert.Equal(t, want, got)
}
