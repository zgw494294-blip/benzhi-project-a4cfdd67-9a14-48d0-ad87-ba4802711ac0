package httpapi

import (
	"bytes"
	"encoding/json"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/workflow"
	"io"
	"net/http"
)

func (s *Server) addSource(w http.ResponseWriter, r *http.Request, id string) {
	body, e := io.ReadAll(r.Body)
	if e != nil {
		respondError(w, e)
		return
	}
	var fields map[string]json.RawMessage
	if e = json.Unmarshal(body, &fields); e != nil {
		respondError(w, e)
		return
	}
	if _, bulk := fields["sources"]; bulk {
		var req workflow.AddSourcesRequest
		if e = json.Unmarshal(body, &req); e != nil {
			respondError(w, e)
			return
		}
		if req.ExpectedVersion == 0 {
			if v, getErr := s.Service.GetBatch(id); getErr == nil {
				req.ExpectedVersion = expected(r, v.Batch.Version)
			}
		}
		out, err := s.Service.AddSources(id, req, req.ExpectedVersion, key(r))
		if err != nil {
			respondError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"batch_id": id, "sources": out})
		return
	}
	var req workflow.AddSourceRequest
	if e = json.NewDecoder(bytes.NewReader(body)).Decode(&req); e != nil {
		respondError(w, e)
		return
	}
	v, e := s.Service.GetBatch(id)
	if e != nil {
		respondError(w, e)
		return
	}
	out, e := s.Service.AddSource(id, req, expected(r, v.Batch.Version), key(r))
	if e != nil {
		respondError(w, e)
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

func (s *Server) policy(w http.ResponseWriter, r *http.Request, id string) {
	out, err := s.Service.PolicyPrecheck(id)
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) status(w http.ResponseWriter, r *http.Request, id string) {
	out, err := s.Service.Status(id)
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) export(w http.ResponseWriter, r *http.Request, id string) {
	out, err := s.Service.Export(id)
	if err != nil {
		respondError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(out)
}

func (s *Server) verify(w http.ResponseWriter, r *http.Request, id string) {
	if err := s.Service.Verify(id); err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"batch_id": id, "valid": true})
}
func (s *Server) review(w http.ResponseWriter, r *http.Request, id, source string) {
	var req workflow.ReviewRequest
	if e := decode(r, &req); e != nil {
		respondError(w, e)
		return
	}
	if req.ExpectedVersion == 0 {
		if v, e := s.Service.GetBatch(id); e == nil {
			req.ExpectedVersion = expected(r, v.Batch.Version)
		}
	}
	out, e := s.Service.Review(id, source, req, key(r))
	if e != nil {
		respondError(w, e)
		return
	}
	writeJSON(w, http.StatusOK, out)
}
func (s *Server) resubmit(w http.ResponseWriter, r *http.Request, id string) {
	var req workflow.ResubmitRequest
	if e := decode(r, &req); e != nil {
		respondError(w, e)
		return
	}
	if req.ExpectedVersion == 0 {
		if v, e := s.Service.GetBatch(id); e == nil {
			req.ExpectedVersion = expected(r, v.Batch.Version)
		}
	}
	out, e := s.Service.Resubmit(id, req, key(r))
	if e != nil {
		respondError(w, e)
		return
	}
	writeJSON(w, http.StatusOK, out)
}
func (s *Server) freeze(w http.ResponseWriter, r *http.Request, id string) {
	var req workflow.ActorRequest
	if r.Body != nil {
		_ = decode(r, &req)
	}
	if req.Actor == "" {
		req.Actor = "发布负责人"
	}
	if req.ExpectedVersion == 0 {
		if v, e := s.Service.GetBatch(id); e == nil {
			req.ExpectedVersion = expected(r, v.Batch.Version)
		}
	}
	out, e := s.Service.Freeze(id, req, key(r))
	if e != nil {
		respondError(w, e)
		return
	}
	writeJSON(w, http.StatusOK, out)
}
func (s *Server) close(w http.ResponseWriter, r *http.Request, id string) {
	var req workflow.ActorRequest
	if r.Body != nil {
		_ = decode(r, &req)
	}
	if req.Actor == "" {
		req.Actor = "发布负责人"
	}
	if req.ExpectedVersion == 0 {
		if v, e := s.Service.GetBatch(id); e == nil {
			req.ExpectedVersion = expected(r, v.Batch.Version)
		}
	}
	out, e := s.Service.Close(id, req, key(r))
	if e != nil {
		respondError(w, e)
		return
	}
	writeJSON(w, http.StatusOK, out)
}
func (s *Server) selfcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	if e := s.Service.SelfCheck(); e != nil {
		respondError(w, e)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
