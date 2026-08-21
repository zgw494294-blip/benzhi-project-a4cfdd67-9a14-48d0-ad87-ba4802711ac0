package httpapi

import (
	"net/http"
	"strings"
)

func requireJSON(r *http.Request) bool {
	return strings.Contains(strings.ToLower(r.Header.Get("Content-Type")), "application/json")
}
func requireKey(r *http.Request) string { return strings.TrimSpace(r.Header.Get("Idempotency-Key")) }
func requireActor(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "匿名操作者"
	}
	return value
}
func noBody(status int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(status) }
}
func validID(id string) bool {
	if len(id) < 3 || len(id) > 128 {
		return false
	}
	for _, r := range id {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_') {
			return false
		}
	}
	return true
}
