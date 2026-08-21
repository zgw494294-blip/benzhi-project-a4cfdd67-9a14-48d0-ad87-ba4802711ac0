package store

import (
	"fmt"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/audit"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
)

type ConsistencyReport struct {
	Batches int      `json:"batches"`
	Sources int      `json:"sources"`
	Reviews int      `json:"reviews"`
	Events  int      `json:"events"`
	Valid   bool     `json:"valid"`
	Errors  []string `json:"errors,omitempty"`
}

func (l *Ledger) CheckConsistency() ConsistencyReport {
	r := ConsistencyReport{}
	_ = l.Read(func(s Snapshot) error {
		r.Batches = len(s.Batches)
		r.Sources = len(s.Sources)
		r.Reviews = len(s.Reviews)
		r.Events = len(s.Events)
		for id, b := range s.Batches {
			sources := []model.SourceRecord{}
			reviews := []model.ReviewTask{}
			for _, x := range s.Sources {
				if x.BatchID == id {
					sources = append(sources, x)
				}
			}
			for _, x := range s.Reviews {
				if x.BatchID == id {
					reviews = append(reviews, x)
				}
			}
			var snap *model.ReleaseSnapshot
			for _, x := range s.Snapshots {
				if x.BatchID == id {
					v := x
					snap = &v
				}
			}
			var cert *model.AuditCertificate
			for _, x := range s.Certificates {
				if x.BatchID == id {
					v := x
					cert = &v
				}
			}
			if err := model.ValidateAggregate(b, sources, reviews, snap, cert); err != nil {
				r.Errors = append(r.Errors, fmt.Sprintf("%s: %v", id, err))
			}
		}
		byBatch := map[string][]model.AuditEvent{}
		for _, ev := range s.Events {
			byBatch[ev.BatchID] = append(byBatch[ev.BatchID], ev)
		}
		for id, events := range byBatch {
			if err := audit.Verify(events); err != nil {
				r.Errors = append(r.Errors, fmt.Sprintf("%s: %v", id, err))
			}
		}
		return nil
	})
	r.Valid = len(r.Errors) == 0
	return r
}
func (l *Ledger) EventRange(batchID string, start, end int64) []model.AuditEvent {
	out := []model.AuditEvent{}
	_ = l.Read(func(s Snapshot) error {
		for _, e := range s.Events {
			if e.BatchID == batchID && (start <= 0 || e.Sequence >= start) && (end <= 0 || e.Sequence <= end) {
				out = append(out, e)
			}
		}
		return nil
	})
	return out
}
