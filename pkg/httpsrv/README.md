# HttpSrv module

This module, pkg or folder contain `helpers`, `server methods` and `middlewares`

```bash
    │   └── middlewares
    │       ├── InjectSecureSpacePath.go
    │       ├── allowacces.go
    │       ├── jwt.go
    │       ├── multitenant.go
    │       ├── o11y.go
    │       └── types.go
```

The `helpers.go` file contain useful methods for management of `request` and `response` in our `handlers`

```go

// M basic type to define a custom map
type M map[string]interface{}

// Bind helper to parse request Body
func Bind(r *http.Request, payload interface{}) error {
	return json.NewDecoder(r.Body).Decode(payload)
}

// ResponseError return error message
func ResponseError(w http.ResponseWriter, code int, msg string) {
	ResponseJSON(w, code, map[string]string{"message": msg})
}

// ResponseJSON write json response format
func ResponseJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
```

In the `httpsrv.go` file are defined the main logic for setting up our http server using ssl or not depend on our context.



### TODO
- Implementar middleware append space
- implementar middleware append owner