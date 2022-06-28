package middlewares

import (
	"context"
	"encoding/json"
	"github.com/piqba/common/pkg/dql"
	"net/http"
)

const (
	CommonCtxKey          = "common"
	AuthorizationHeader   = "Authorization"
	WorkSpaceXHeader      = "x-workspace"
	SpaceMultipleProperty = "multiple"
	WorkSpaceXHeaderEmpty = ""
	LanguageHeader        = "accepted-language"
)

type (
	// CommonCtx define all objects in the context custom
	CommonCtx struct {
		RequestSearchCtx RequestSearchCtx `json:"request_search_ctx"`
		ClaimCtx         ClaimCtx         `json:"claim_ctx"`
	}
	// RequestSearchCtx ...
	RequestSearchCtx struct {
		Size    int         `json:"size"`
		Sort    string      `json:"sort"`
		Page    int         `json:"page"`
		Search  string      `json:"search"`
		Filters dql.Filters `json:"filters"`
	}
	// ClaimCtx return an struct with all value from token claims
	ClaimCtx struct {
		permissions []string
		email       string
		Space       Space `json:"space"`
		userId      string
	}
	// Space ...
	Space struct {
		Id              string   `json:"_id"`
		Identifier      string   `json:"identifier"`
		IsPublicSpace   bool     `json:"isPublicSpace"`
		IsRootSpace     bool     `json:"isRootSpace"`
		SecureSpacePath []string `json:"secureSpacePath"`
	}
	//CacheCtx ...
	CacheCtx struct {
		keyPrefix string
	}
)

func NewRequestSearchDefault() *RequestSearchCtx {
	return &RequestSearchCtx{
		Size:   20,
		Sort:   "desc",
		Page:   1,
		Search: "",
	}
}

// DefaultPage ...
func (rctx RequestSearchCtx) DefaultPage(value, fallback int) int {
	if value != 0 {
		return value
	}
	return fallback
}

// DefaultSize ...
func (rctx RequestSearchCtx) DefaultSize(value, fallback int) int {
	if value != 0 {
		return value
	}
	return fallback
}

// DefaultSort ...
func (rctx RequestSearchCtx) DefaultSort(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

// GetFilterJsonString ...
func (rctx RequestSearchCtx) GetFilterJsonString() string {
	bytes, _ := json.Marshal(rctx.Filters)
	return string(bytes)
}

// NewCommonContext returns a new Context that carries value u.
func NewCommonContext(ctx context.Context, key string, c *CommonCtx) context.Context {
	return context.WithValue(ctx, key, c)
}

// FromCommonContext returns the User value stored in ctx, if any.
func FromCommonContext(ctx context.Context, key string) (*CommonCtx, bool) {
	c, ok := ctx.Value(key).(*CommonCtx)
	return c, ok
}

// Middleware type is a slice of standard middleware handlers with methods
// to compose middleware chains and http.Handler's.
type Middleware func(http.Handler) http.Handler

// Chain applies middlewares to a APIGatewayHandlerFunc
func Chain(f http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}
