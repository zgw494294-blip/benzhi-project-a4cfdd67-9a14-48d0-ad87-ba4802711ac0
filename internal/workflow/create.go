package workflow

import (
	"context"
	"encoding/json"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/store"
	"sync/atomic"
	"time"
)

var idCounter uint64

func (s *Service) CreateBatch(r CreateBatchRequest, key string) (model.DatasetBatch, error) {
	return s.CreateBatchContext(context.Background(), r, key)
}

func (s *Service) CreateBatchContext(ctx context.Context, r CreateBatchRequest, key string) (model.DatasetBatch, error) {
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
	err := s.Ledger.UpdateContext(ctx, func(st *store.Snapshot) error {
		st.Batches[out.BatchID] = out
		s.appendEvent(st, out.BatchID, model.NewEvent("batch_created", r.Steward, map[string]any{"title": r.Title}))
		if key != "" {
			data, marshalErr := json.Marshal(out)
			if marshalErr != nil {
				return marshalErr
			}
			st.Idempotency[store.IdempotencyKey("create", key)] = data
		}
		return nil
	})
	if err != nil {
		return model.DatasetBatch{}, err
	}
	return out, nil
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
