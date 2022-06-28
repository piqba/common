package middlewares

import (
	"encoding/json"
	"github.com/piqba/common/pkg/httpsrv"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAllowAccess_UnknownACL(t *testing.T) {
	assertions := assert.New(t)
	want := httpsrv.M{
		"message": "error: you don`t have permissions to view this content",
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	ctx := NewCommonContext(req.Context(), CommonCtxKey, &CommonCtx{
		RequestSearchCtx: RequestSearchCtx{},
		ClaimCtx:         claimTest,
	})

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	h := Chain(
		nextHandler,
		AllowAccess("bad"),
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
func TestAllowAccess_OK(t *testing.T) {
	assertions := assert.New(t)
	want := httpsrv.M{}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	ctx := NewCommonContext(req.Context(), CommonCtxKey, &CommonCtx{
		RequestSearchCtx: RequestSearchCtx{},
		ClaimCtx:         claimTest,
	})

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	h := Chain(
		nextHandler,
		AllowAccess(FullAcl, AdminAcl, GuestAcl),
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
func TestAllowAccess_FullACL(t *testing.T) {
	assertions := assert.New(t)
	want := httpsrv.M{}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	ctx := NewCommonContext(req.Context(), CommonCtxKey, &CommonCtx{
		RequestSearchCtx: RequestSearchCtx{},
		ClaimCtx:         claimTest,
	})

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	h := Chain(
		nextHandler,
		AllowAccess(FullAcl),
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
