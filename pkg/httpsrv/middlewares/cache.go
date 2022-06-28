package middlewares

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/piqba/common/pkg/cache/factory"
	"github.com/piqba/common/pkg/httpsrv"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	keyStrReplacer = strings.NewReplacer("$", "")
)

const (
	SpaceKW  KeyWord = "space"  // get space from context
	UserKW   KeyWord = "user"   // get space from context
	LangKW   KeyWord = "lang"   // header accepted-language
	PageKW   KeyWord = "page"   // get from body object request
	FilterKW KeyWord = "filter" // get from body object request
	SearchKW KeyWord = "search" // get from search object request
)

type (
	KeyWord         string
	CacheMiddleware struct {
		key   string
		cache factory.ICache
		ttl   time.Duration
	}
	CacheResponseWriter struct {
		http.ResponseWriter
		buf *bytes.Buffer
		ctx context.Context
		chm CacheMiddleware
	}
	CacheResponseData struct {
		Data       []interface{} `json:"data"`
		SimpleData interface{}   `json:"simple_data"`
		Source     string        `json:"source"`
	}
)

func (mrw *CacheResponseWriter) Write(p []byte) (int, error) {

	//TODO: For check when the backend send [] or zero default value
	fmt.Println(len(p))

	err := mrw.chm.cache.Set(mrw.ctx, []byte(mrw.chm.key), p, mrw.chm.ttl)
	if err != nil {
		return 0, err
	}
	return mrw.buf.Write(p)

}

// hashGenerator return hash 256
func hashGenerator(prefixKey string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(prefixKey))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	return hash
}

// buildPrefixKey return key string from request and prefixKeyOptions parameters
func buildPrefixKey(r *http.Request, prefixKeyOptions []interface{}) string {
	var keys []string
	for _, opt := range prefixKeyOptions {
		switch kopt := opt.(type) {
		case string:
			switch {
			case strings.Contains(kopt, "$"):
				keys = append(
					keys,
					keyStrReplacer.Replace(kopt),
				)
			case strings.Contains(kopt, "params"):
				param := strings.Split(kopt, "params.")[1]
				keys = append(
					keys,
					hashGenerator(
						chi.URLParam(r, param),
					),
				)
			case strings.Contains(kopt, "query"):
				// TODO: r.URL.Query()
			}
		case KeyWord:
			keys = append(
				keys,
				hashGenerator(
					extractValueFromCtxByKW(r, kopt),
				),
			)
		}
	}
	return strings.Join(keys, ":")
}

// extractValueFromCtxByKW ...
func extractValueFromCtxByKW(r *http.Request, word KeyWord) string {
	var key string
	common, ok := FromCommonContext(r.Context(), CommonCtxKey)
	if ok {
		switch word {
		case SpaceKW:
			key = common.ClaimCtx.Space.Id
		case UserKW:
			key = common.ClaimCtx.userId
		case LangKW:
			key = r.Header.Get(LanguageHeader)
		case SearchKW:
			key = common.RequestSearchCtx.Search
		case PageKW:
			key = strconv.Itoa(common.RequestSearchCtx.Page)
		case FilterKW:
			key = common.RequestSearchCtx.GetFilterJsonString()
		}
	}
	return key
}

// CacheDFL ...
func CacheDFL(
	cache factory.ICache,
	prefixKeyOptions ...interface{},
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		cacheFn := func(w http.ResponseWriter, r *http.Request) {
			key := buildPrefixKey(r, prefixKeyOptions)
			b, err := cache.Get(r.Context(), []byte(key))
			switch {
			case err == redis.Nil:
				mrw := &CacheResponseWriter{
					ResponseWriter: w,
					buf:            &bytes.Buffer{},
					ctx:            r.Context(),
					chm: CacheMiddleware{
						key:   key,
						cache: cache,
						ttl:   time.Hour,
					},
				}
				next.ServeHTTP(mrw, r.WithContext(r.Context()))
				if _, err := io.Copy(w, mrw.buf); err != nil {
					httpsrv.ResponseError(w, http.StatusInternalServerError, "error to save data from Cache")
					return
				}
				return
			case err != nil:
				httpsrv.ResponseError(w, http.StatusInternalServerError, "error to save data from Cache")
				return
			default:
				var responseData CacheResponseData
				err = json.Unmarshal(b, &responseData.Data)
				if err != nil {
					err2 := json.Unmarshal(b, &responseData.SimpleData)
					if err2 != nil {
						httpsrv.ResponseError(w, http.StatusInternalServerError, "error to Unmarshal data from Cache")
						return
					}
				}
				responseData.Source = "cache"
				httpsrv.ResponseJSON(w, http.StatusOK, responseData)
			}
		}
		return http.HandlerFunc(cacheFn)
	}
}

// InvalidCacheDFL ...
func InvalidCacheDFL(
	cache factory.ICache,
	prefixKeyOptions ...interface{},
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		cacheFn := func(w http.ResponseWriter, r *http.Request) {
			switch strings.ToLower(r.Method) {
			case "patch", "put", "delete", "post":
				// Invalid From Redis
				key := buildPrefixKey(r, prefixKeyOptions)
				_, err := cache.Invalidate(r.Context(), key)
				if err != nil {
					httpsrv.ResponseError(w, http.StatusInternalServerError, "error to invalidated data from Cache")
					return
				}
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(cacheFn)
	}
}

// CacheBasicKey ...
func CacheBasicKey(
	cache factory.ICache, prefixKey string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		cacheFn := func(w http.ResponseWriter, r *http.Request) {
			// prepare key identifier
			key := fmt.Sprintf("%s:%s", prefixKey, r.URL.Path)

			// Get From Redis
			b, err := cache.Get(r.Context(), []byte(key))
			if err == redis.Nil {
				log.Println("key not found:", err)
			} else if err != nil {
				httpsrv.ResponseError(w, http.StatusInternalServerError, "error to save data from Cache")
				return
			}
			// validate redis response
			if len(b) == 0 {
				mrw := &CacheResponseWriter{
					ResponseWriter: w,
					buf:            &bytes.Buffer{},
					ctx:            r.Context(),
					chm: CacheMiddleware{
						key:   key,
						cache: cache,
						ttl:   time.Hour,
					},
				}
				next.ServeHTTP(mrw, r.WithContext(r.Context()))
				if _, err := io.Copy(w, mrw.buf); err != nil {
					httpsrv.ResponseError(w, http.StatusInternalServerError, "error to save data from Cache")
					return
				}
				return
			}
			var responseData CacheResponseData
			err = json.Unmarshal(b, &responseData.Data)
			if err != nil {
				err2 := json.Unmarshal(b, &responseData.SimpleData)
				if err2 != nil {
					httpsrv.ResponseError(w, http.StatusInternalServerError, "error to Unmarshal data from Cache")
					return
				}
			}
			responseData.Source = "cache"
			httpsrv.ResponseJSON(w, http.StatusOK, responseData)

		}
		return http.HandlerFunc(cacheFn)
	}
}

// InvalidCacheBasicKey ...
func InvalidCacheBasicKey(
	cache factory.ICache, prefixKey string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		cacheFn := func(w http.ResponseWriter, r *http.Request) {
			switch strings.ToLower(r.Method) {
			case "patch", "put", "delete", "post":
				// Invalid From Redis
				_, err := cache.Invalidate(r.Context(), prefixKey)
				if err != nil {
					httpsrv.ResponseError(w, http.StatusInternalServerError, "error to invalidated data from Cache")
					return
				}
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(cacheFn)
	}
}
