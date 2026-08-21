package workflow

import (
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/store"
)

func (s *Service) Freeze(batchID string, r ActorRequest, key string) (model.ReleaseSnapshot, error) {
	var out model.ReleaseSnapshot
	if key != "" {
		if ok, err := s.Ledger.GetIdempotent("freeze:"+batchID, key, &out); ok || err != nil {
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
		sources := sourceRecords(st, batchID)
		if e = model.EnsureCanFreeze(b, sources); e != nil {
			return e
		}
		now := model.Now()
		out = model.ReleaseSnapshot{SnapshotID: newID("snapshot"), BatchID: batchID, ManifestHash: model.ManifestHash(b, sources), ApprovedScope: b.ReleaseScope, SourceCount: len(sources), FrozenAt: now, FrozenBy: r.Actor}
		st.Snapshots[out.SnapshotID] = out
		b.Status = model.StatusFrozen
		b.Version++
		b.UpdatedAt = now
		st.Batches[batchID] = b
		s.appendEvent(st, batchID, model.NewEvent("batch_frozen", r.Actor, map[string]any{"snapshot_id": out.SnapshotID, "manifest_hash": out.ManifestHash}))
		return nil
	})
	if err == nil && key != "" {
		err = s.Ledger.PutIdempotent("freeze:"+batchID, key, out)
	}
	return out, err
}
