package workflow

import (
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/store"
	"sync/atomic"
	"time"
)

var idCounter uint64

func (s *Service) CreateBatch(r CreateBatchRequest, key string) (model.DatasetBatch, error) {
	var out model.DatasetBatch
	if key != "" {
		if ok, err := s.Ledger.GetIdempotent("create", key, &out); ok || err != nil {
			return out, err
		}
	}
	if err := model.ValidateBatchInput(r.Title, r.Steward, r.PolicyVersion, r.ReleaseScope); err != nil {
		return out, err
	}
	now := model.Now()
	out = model.DatasetBatch{BatchID: newID("batch"), Title: r.Title, Steward: r.Steward, PolicyVersion: r.PolicyVersion, ReleaseScope: r.ReleaseScope, Status: model.StatusDraft, Version: 1, CreatedAt: now, UpdatedAt: now}
	err := s.Ledger.Update(func(st *store.Snapshot) error {
		st.Batches[out.BatchID] = out
		s.appendEvent(st, out.BatchID, model.NewEvent("batch_created", r.Steward, map[string]any{"title": r.Title}))
		return nil
	})
	if err == nil && key != "" {
		err = s.Ledger.PutIdempotent("create", key, out)
	}
	return out, err
}
func newID(prefix string) string {
	return prefix + "-" + model.HashJSON(struct {
		T string
		N any
	}{prefix, struct {
		Time string
		Seq  uint64
	}{model.Now().Format(time.RFC3339Nano), atomic.AddUint64(&idCounter, 1)}})[:20]
}
