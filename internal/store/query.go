package store

import "github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"

func (l *Ledger) Batch(id string) (model.DatasetBatch, error) {
	var out model.DatasetBatch
	err := l.Read(func(s Snapshot) error {
		v, ok := s.Batches[id]
		if !ok {
			return model.NotFound("batch", id)
		}
		out = v
		return nil
	})
	return out, err
}
func (l *Ledger) Sources(batchID string) []model.SourceRecord {
	out := []model.SourceRecord{}
	_ = l.Read(func(s Snapshot) error {
		for _, v := range s.Sources {
			if v.BatchID == batchID {
				out = append(out, v)
			}
		}
		model.SortSources(out)
		return nil
	})
	return out
}
func (l *Ledger) Reviews(batchID string) []model.ReviewTask {
	out := []model.ReviewTask{}
	_ = l.Read(func(s Snapshot) error {
		for _, v := range s.Reviews {
			if v.BatchID == batchID {
				out = append(out, v)
			}
		}
		model.SortReviews(out)
		return nil
	})
	return out
}
func (l *Ledger) Snapshot(batchID string) (model.ReleaseSnapshot, bool) {
	var out model.ReleaseSnapshot
	ok := false
	_ = l.Read(func(s Snapshot) error {
		for _, v := range s.Snapshots {
			if v.BatchID == batchID {
				out = v
				ok = true
				break
			}
		}
		return nil
	})
	return out, ok
}
func (l *Ledger) Certificate(batchID string) (model.AuditCertificate, bool) {
	var out model.AuditCertificate
	ok := false
	_ = l.Read(func(s Snapshot) error {
		for _, v := range s.Certificates {
			if v.BatchID == batchID {
				out = v
				ok = true
				break
			}
		}
		return nil
	})
	return out, ok
}
func (l *Ledger) Count() (int, int, int) {
	b, s, e := 0, 0, 0
	_ = l.Read(func(x Snapshot) error { b = len(x.Batches); s = len(x.Sources); e = len(x.Events); return nil })
	return b, s, e
}
