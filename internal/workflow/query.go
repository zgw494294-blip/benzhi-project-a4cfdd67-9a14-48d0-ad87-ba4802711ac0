package workflow

import (
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/store"
	"sort"
)

func (s *Service) GetBatch(id string) (model.BatchView, error) {
	var out model.BatchView
	err := s.Ledger.Read(func(st store.Snapshot) error {
		b, e := getBatch(&st, id)
		if e != nil {
			return e
		}
		out.Batch = b
		for _, x := range st.Sources {
			if x.BatchID == id {
				out.Sources = append(out.Sources, x)
			}
		}
		for _, x := range st.Reviews {
			if x.BatchID == id {
				out.Reviews = append(out.Reviews, x)
			}
		}
		for _, x := range st.Snapshots {
			if x.BatchID == id {
				v := x
				out.Snapshot = &v
			}
		}
		for _, x := range st.Certificates {
			if x.BatchID == id {
				v := x
				out.Certificate = &v
			}
		}
		out.Audit = filterEvents(st.Events, id)
		model.SortSources(out.Sources)
		model.SortReviews(out.Reviews)
		sort.Slice(out.Audit, func(i, j int) bool { return out.Audit[i].Sequence < out.Audit[j].Sequence })
		return nil
	})
	return out, err
}
func (s *Service) GetAudit(id string) ([]model.AuditEvent, error) {
	var out []model.AuditEvent
	err := s.Ledger.Read(func(st store.Snapshot) error {
		if _, e := getBatch(&st, id); e != nil {
			return e
		}
		out = filterEvents(st.Events, id)
		return nil
	})
	return out, err
}
