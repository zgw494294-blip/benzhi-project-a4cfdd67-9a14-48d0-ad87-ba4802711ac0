package workflow

import (
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/store"
)

func (s *Service) Review(batchID, sourceID string, r ReviewRequest, key string) (model.ReviewTask, error) {
	var out model.ReviewTask
	if key != "" {
		if ok, err := s.Ledger.GetIdempotent("review:"+batchID+":"+sourceID, key, &out); ok || err != nil {
			return out, err
		}
	}
	if err := model.ValidateReviewInput(r.Reviewer, r.Decision, r.EvidenceRefs, r.Issues); err != nil {
		return out, err
	}
	err := s.Ledger.Update(func(st *store.Snapshot) error {
		b, e := getBatch(st, batchID)
		if e != nil {
			return e
		}
		if e = checkVersion(b, r.ExpectedVersion); e != nil {
			return e
		}
		src, ok := st.Sources[sourceID]
		if !ok || src.BatchID != batchID {
			return model.NotFound("source", sourceID)
		}
		if e = model.EnsureCanReview(b, sourceRecords(st, batchID), sourceID); e != nil {
			return e
		}
		now := model.Now()
		revision := 1
		for _, prior := range st.Reviews {
			if prior.BatchID == batchID && prior.SourceID == sourceID && prior.Revision >= revision {
				revision = prior.Revision + 1
			}
		}
		evidence := append([]string(nil), r.EvidenceRefs...)
		if revision > 1 {
			evidence = MergeEvidence(src.EvidenceRefs, evidence)
			for _, prior := range st.Reviews {
				if prior.BatchID == batchID && prior.SourceID == sourceID {
					evidence = MergeEvidence(evidence, prior.EvidenceRefs)
				}
			}
		}
		out = model.ReviewTask{ReviewID: newID("review"), BatchID: batchID, SourceID: sourceID, Reviewer: r.Reviewer, Decision: r.Decision, Issues: append([]string(nil), r.Issues...), EvidenceRefs: evidence, ReviewedAt: now, Revision: revision}
		st.Reviews[out.ReviewID] = out
		switch r.Decision {
		case model.DecisionApprove:
			src.State = model.SourceApproved
		case model.DecisionCorrection, model.DecisionReject:
			src.State = model.SourceNeedsCorrection
			b.Status = model.StatusCorrection
		}
		src.EvidenceRefs = MergeEvidence(src.EvidenceRefs, evidence)
		st.Sources[sourceID] = src
		b.Version++
		b.UpdatedAt = now
		st.Batches[batchID] = b
		s.appendEvent(st, batchID, model.NewEvent("source_reviewed", r.Reviewer, map[string]any{"source_id": sourceID, "decision": r.Decision, "issues": r.Issues}))
		return nil
	})
	if err == nil && key != "" {
		err = s.Ledger.PutIdempotent("review:"+batchID+":"+sourceID, key, out)
	}
	return out, err
}
func sourceRecords(st *store.Snapshot, batchID string) []model.SourceRecord {
	out := []model.SourceRecord{}
	for _, x := range st.Sources {
		if x.BatchID == batchID {
			out = append(out, x)
		}
	}
	return out
}
