package middlewares

import (
	"encoding/json"
	"github.com/piqba/common/pkg/httpsrv"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	secretOK  = "development"
	secretBad = "bad"

	tokenTestOK      = "Gincw8IGS8sxPAaMUESUwEo"
	tokenTestExpired = "enu8"
	tokenTestInvalid = "2aMUESUwEo"
)

func TestVerifyJWT_BadHeader(t *testing.T) {
	assertions := assert.New(t)

	want := httpsrv.M{
		"message": "error: Format is Authorization: Bearer [token]",
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	req.Header.Set(AuthorizationHeader, "as")
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(w)
	})
	h1 := VerifyJWT("")(nextHandler)
	h1.ServeHTTP(w, req)

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
func TestVerifyJWT_TokenValidator_Invalid_Secret(t *testing.T) {
	assertions := assert.New(t)
	want := httpsrv.M{
		"message": "couldn't handle this token",
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	req.Header.Set(AuthorizationHeader, "Bearer "+tokenTestOK)
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(w)
	})
	h1 := VerifyJWT(secretBad)(nextHandler)
	h1.ServeHTTP(w, req)

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
func TestVerifyJWT_TokenValidator_Expired_TOKEN(t *testing.T) {
	assertions := assert.New(t)
	want := httpsrv.M{
		"message": "error: Token is either expired or not active yet",
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	req.Header.Set(AuthorizationHeader, "Bearer "+tokenTestExpired)
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(w)
	})
	h1 := VerifyJWT("")(nextHandler)
	h1.ServeHTTP(w, req)

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
func TestVerifyJWT_TokenValidator_Invalid_TOKEN(t *testing.T) {
	assertions := assert.New(t)
	want := httpsrv.M{
		"message": "error: That's not even a token",
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	req.Header.Set(AuthorizationHeader, "Bearer "+tokenTestInvalid)
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(w)
	})
	h1 := VerifyJWT("")(nextHandler)
	h1.ServeHTTP(w, req)

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
func TestVerifyJWT_OK(t *testing.T) {
	assertions := assert.New(t)
	want := httpsrv.M{}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	req.Header.Set(AuthorizationHeader, "Bearer "+tokenTestOK)
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		common, _ := FromCommonContext(r.Context(), CommonCtxKey)
		assert.Equal(t, common.ClaimCtx.Space.Id, "idSpace")
	})
	h1 := VerifyJWT(secretOK)(nextHandler)
	h1.ServeHTTP(w, req)

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
