package middlewares

import (
	"github.com/go-redis/redis/v8"
	"github.com/piqba/common/pkg/cache"
	"github.com/piqba/common/pkg/cache/factory"
	"github.com/piqba/common/pkg/httpsrv"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCacheBasicKey_OK(t *testing.T) {
	//assertions := assert.New(t)
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cch, _ := factory.GetCacheFactory(cache.Redis, cache.CacheOptions{
		RdbClient: rdb,
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	ctx := NewCommonContext(req.Context(), CommonCtxKey, &CommonCtx{
		RequestSearchCtx: RequestSearchCtx{},
		ClaimCtx:         claimTest,
	})

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpsrv.ResponseJSON(w, http.StatusOK, httpsrv.M{"data": "ok"})
	})

	h := Chain(
		nextHandler,
		CacheBasicKey(cch, "testKey"),
	)
	h.ServeHTTP(w, req.WithContext(ctx))

	if w.Code != http.StatusOK {
		t.Errorf("Status code spected %d, getter %d", http.StatusOK, w.Code)
	}
}
