package workflow

import (
	"fmt"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/audit"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/store"
)

type Service struct{ Ledger *store.Ledger }

func New(l *store.Ledger) *Service { return &Service{Ledger: l} }
func (s *Service) appendEvent(st *store.Snapshot, batchID string, ev model.Event) {
	chain := audit.New(filterEvents(st.Events, batchID))
	chain.Append(batchID, ev)
	others := make([]model.AuditEvent, 0, len(st.Events)+1)
	for _, existing := range st.Events {
		if existing.BatchID != batchID {
			others = append(others, existing)
		}
	}
	st.Events = append(others, chain.Events()...)
}
func getBatch(st *store.Snapshot, id string) (model.DatasetBatch, error) {
	b, ok := st.Batches[id]
	if !ok {
		return b, model.NotFound("batch", id)
	}
	return b, nil
}
func checkVersion(b model.DatasetBatch, expected int64) error {
	if expected > 0 && b.Version != expected {
		return model.Err("version_mismatch", fmt.Sprintf("版本不匹配，当前版本为 %d", b.Version))
	}
	return nil
}
