// author kenriortega

package httpsrv

import (
	"encoding/json"
	"net/http"
)

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
