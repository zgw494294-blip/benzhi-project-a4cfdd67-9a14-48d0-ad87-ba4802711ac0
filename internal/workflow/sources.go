package workflow

import (
	"encoding/json"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/store"
	"sort"
	"strings"
)

func (s *Service) AddSource(batchID string, r AddSourceRequest, expected int64, key string) (model.SourceRecord, error) {
	var out model.SourceRecord
	if key != "" {
		if ok, err := s.Ledger.GetIdempotent("source:"+batchID, key, &out); ok || err != nil {
			return out, err
		}
	}
	if err := model.ValidateSourceInput(r.Origin, r.Description, r.ConsentRef, r.Sensitivity, r.ContentChecksum); err != nil {
		return out, err
	}
	err := s.Ledger.Update(func(st *store.Snapshot) error {
		b, e := getBatch(st, batchID)
		if e != nil {
			return e
		}
		if e = checkVersion(b, expected); e != nil {
			return e
		}
		if e = model.EnsureCanAddSource(b); e != nil {
			return e
		}
		now := model.Now()
		out = model.SourceRecord{SourceID: newID("source"), BatchID: batchID, Origin: r.Origin, Description: r.Description, ConsentRef: r.ConsentRef, Sensitivity: r.Sensitivity, ContentChecksum: r.ContentChecksum, EvidenceRefs: []string{r.ConsentRef}, State: model.SourcePending, CreatedAt: now}
		st.Sources[out.SourceID] = out
		b.Version++
		b.UpdatedAt = now
		if b.Status == model.StatusDraft {
			b.Status = model.StatusReviewing
		}
		st.Batches[batchID] = b
		s.appendEvent(st, batchID, model.NewEvent("source_added", b.Steward, map[string]any{"source_id": out.SourceID}))
		return nil
	})
	if err == nil && key != "" {
		err = s.Ledger.PutIdempotent("source:"+batchID, key, out)
	}
	return out, err
}

// AddSources validates and persists a complete source registration request atomically.
func (s *Service) AddSources(batchID string, r AddSourcesRequest, expected int64, key string) ([]model.SourceRecord, error) {
	var out []model.SourceRecord
	if key != "" {
		if ok, err := s.Ledger.GetIdempotent("sources-batch:"+batchID, key, &out); ok || err != nil {
			return out, err
		}
	}
	if len(r.Sources) == 0 {
		return nil, model.Invalid("sources 不能为空")
	}
	origins := map[string]bool{}
	checksums := map[string]bool{}
	for _, src := range r.Sources {
		if err := model.ValidateSourceInput(src.Origin, src.Description, src.ConsentRef, src.Sensitivity, src.ContentChecksum); err != nil {
			return nil, err
		}
		origin := strings.TrimSpace(src.Origin)
		checksum := strings.TrimSpace(src.ContentChecksum)
		if origins[origin] {
			return nil, model.Conflict("请求中存在重复 origin")
		}
		if checksums[checksum] {
			return nil, model.Conflict("请求中存在重复 content_checksum")
		}
		origins[origin] = true
		checksums[checksum] = true
	}
	err := s.Ledger.Update(func(st *store.Snapshot) error {
		b, e := getBatch(st, batchID)
		if e != nil {
			return e
		}
		if e = checkVersion(b, expected); e != nil {
			return e
		}
		if e = model.EnsureCanAddSource(b); e != nil {
			return e
		}
		for _, existing := range st.Sources {
			if existing.BatchID != batchID {
				continue
			}
			if origins[strings.TrimSpace(existing.Origin)] {
				return model.Conflict("批次中存在重复 origin")
			}
			if checksums[strings.TrimSpace(existing.ContentChecksum)] {
				return model.Conflict("批次中存在重复 content_checksum")
			}
		}
		now := model.Now()
		out = make([]model.SourceRecord, 0, len(r.Sources))
		for _, req := range r.Sources {
			src := model.SourceRecord{SourceID: newID("source"), BatchID: batchID, Origin: req.Origin, Description: req.Description, ConsentRef: req.ConsentRef, Sensitivity: req.Sensitivity, ContentChecksum: req.ContentChecksum, EvidenceRefs: []string{req.ConsentRef}, State: model.SourcePending, CreatedAt: now}
			st.Sources[src.SourceID] = src
			out = append(out, src)
			s.appendEvent(st, batchID, model.NewEvent("source_added", b.Steward, map[string]any{"source_id": src.SourceID}))
		}
		b.Version++
		b.UpdatedAt = now
		if b.Status == model.StatusDraft {
			b.Status = model.StatusReviewing
		}
		st.Batches[batchID] = b
		if key != "" {
			data, marshalErr := json.Marshal(out)
			if marshalErr != nil {
				return marshalErr
			}
			st.Idempotency[store.IdempotencyKey("sources-batch:"+batchID, key)] = data
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool { return out[i].SourceID < out[j].SourceID })
	return out, nil
}
