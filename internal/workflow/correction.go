package workflow

import (
	"encoding/json"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/store"
	"sort"
	"strings"
)

func (s *Service) Resubmit(batchID string, r ResubmitRequest, key string) (model.DatasetBatch, error) {
	var out model.DatasetBatch
	if key != "" {
		if ok, err := s.Ledger.GetIdempotent("resubmit:"+batchID, key, &out); ok || err != nil {
			return out, err
		}
	}
	err := s.Ledger.Update(func(st *store.Snapshot) error {
		b, e := getBatch(st, batchID)
		if e != nil {
			return e
		}
		if e = checkVersion(b, r.ExpectedVersion); e != nil {
			return e
		}
		if b.Status != model.StatusCorrection {
			return model.Conflict("当前状态不需要补证")
		}
		corrections := map[string]bool{}
		for _, src := range sourceRecords(st, batchID) {
			if src.State == model.SourceNeedsCorrection || src.State == model.SourceRejected {
				corrections[src.SourceID] = true
			}
		}
		if len(corrections) == 0 || len(r.EvidenceRefs) != len(corrections) {
			return model.Invalid("evidence_refs 必须精确覆盖待补证来源")
		}
		ids := make([]string, 0, len(corrections))
		for id := range corrections {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		for id, refs := range r.EvidenceRefs {
			if !corrections[id] {
				if src, ok := st.Sources[id]; !ok || src.BatchID != batchID {
					return model.NotFound("source", id)
				}
				return model.Conflict("补证来源必须是当前待补证来源")
			}
			if len(refs) == 0 {
				return model.Invalid("evidence_refs 不能为空")
			}
			for _, ref := range refs {
				if strings.TrimSpace(ref) == "" {
					return model.Invalid("evidence_refs 不能包含空项")
				}
			}
		}
		for _, id := range ids {
			refs := r.EvidenceRefs[id]
			src := st.Sources[id]
			history := src.EvidenceRefs
			for _, prior := range st.Reviews {
				if prior.BatchID == batchID && prior.SourceID == id {
					history = MergeEvidence(history, prior.EvidenceRefs)
				}
			}
			merged := MergeEvidence(history, refs)
			src.EvidenceRefs = merged
			src.ConsentRef = refs[0]
			src.State = model.SourcePending
			st.Sources[id] = src
		}
		b.Status = model.StatusReviewing
		b.Version++
		b.UpdatedAt = model.Now()
		st.Batches[batchID] = b
		out = b
		s.appendEvent(st, batchID, model.NewEvent("correction_resubmitted", b.Steward, map[string]any{"sources": ids}))
		if key != "" {
			data, marshalErr := json.Marshal(out)
			if marshalErr != nil {
				return marshalErr
			}
			st.Idempotency[store.IdempotencyKey("resubmit:"+batchID, key)] = data
		}
		return nil
	})
	return out, err
}
