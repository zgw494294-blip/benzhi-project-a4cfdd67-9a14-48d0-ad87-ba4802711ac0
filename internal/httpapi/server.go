package httpapi

import (
	"encoding/json"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/workflow"
	"net/http"
	"strings"
)

type Server struct {
	Service *workflow.Service
	Mux     *http.ServeMux
}

func New(s *workflow.Service) *Server {
	x := &Server{Service: s, Mux: http.NewServeMux()}
	x.register()
	return x
}
func (s *Server) register() {
	s.Mux.HandleFunc("/v1/batches", s.batches)
	s.Mux.HandleFunc("/v1/batches/", s.batchSub)
	s.Mux.HandleFunc("/selfcheck", s.selfcheck)
}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) { s.Mux.ServeHTTP(w, r) }
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func decode(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return r.Context().Err()
}
func key(r *http.Request) string          { return strings.TrimSpace(r.Header.Get("Idempotency-Key")) }
