package httpapi

import (
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/workflow"
	"net/http"
	"strconv"
	"strings"
)

func expected(r *http.Request, body int64) int64 {
	if v := strings.TrimSpace(r.Header.Get("If-Match")); v != "" {
		n, _ := strconv.ParseInt(v, 10, 64)
		return n
	}
	return body
}
func (s *Server) batches(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req workflow.CreateBatchRequest
		if err := decode(r, &req); err != nil {
			respondError(w, err)
			return
		}
		if err := r.Context().Err(); err != nil {
			respondError(w, err)
			return
		}
		v, e := s.Service.CreateBatchContext(r.Context(), req, key(r))
		if e != nil {
			respondError(w, e)
			return
		}
		writeJSON(w, http.StatusCreated, v)
		return
	}
	http.NotFound(w, r)
}
func (s *Server) batchSub(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1/batches/")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 1 && r.Method == http.MethodGet {
		v, e := s.Service.GetBatch(parts[0])
		if e != nil {
			respondError(w, e)
			return
		}
		writeJSON(w, http.StatusOK, v)
		return
	}
	if len(parts) == 2 && parts[1] == "sources" && r.Method == http.MethodPost {
		s.addSource(w, r, parts[0])
		return
	}
	if len(parts) == 3 && parts[1] == "sources" && (parts[2] == "batch" || parts[2] == "bulk") && r.Method == http.MethodPost {
		s.addSource(w, r, parts[0])
		return
	}
	if len(parts) == 3 && r.Method == http.MethodGet {
		if (parts[1] == "policy" && parts[2] == "precheck") || (parts[1] == "release" && parts[2] == "precheck") {
			s.policy(w, r, parts[0])
			return
		}
		if parts[1] == "certificate" && parts[2] == "verify" {
			s.verify(w, r, parts[0])
			return
		}
	}
	if len(parts) == 2 && parts[1] == "freeze" && r.Method == http.MethodPost {
		s.freeze(w, r, parts[0])
		return
	}
	if len(parts) == 2 && parts[1] == "close" && r.Method == http.MethodPost {
		s.close(w, r, parts[0])
		return
	}
	if len(parts) == 2 && parts[1] == "resubmit" && r.Method == http.MethodPost {
		s.resubmit(w, r, parts[0])
		return
	}
	if len(parts) == 3 && parts[1] == "reviews" && r.Method == http.MethodPost {
		s.review(w, r, parts[0], parts[2])
		return
	}
	if len(parts) == 2 && parts[1] == "audit" && r.Method == http.MethodGet {
		v, e := s.Service.GetAudit(parts[0])
		if e != nil {
			respondError(w, e)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"batch_id": parts[0], "audit": v})
		return
	}
	if len(parts) == 2 && r.Method == http.MethodGet {
		switch parts[1] {
		case "policy", "precheck", "policy-precheck", "preflight", "release-scope":
			s.policy(w, r, parts[0])
			return
		case "status":
			s.status(w, r, parts[0])
			return
		case "export", "artifact":
			s.export(w, r, parts[0])
			return
		case "verify", "certificate", "certificate-verify":
			s.verify(w, r, parts[0])
			return
		}
	}
	http.NotFound(w, r)
}
