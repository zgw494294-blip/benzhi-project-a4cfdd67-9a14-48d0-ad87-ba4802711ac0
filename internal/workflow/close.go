package workflow

import (
	"encoding/json"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/audit"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/store"
)

func (s *Service) Close(batchID string, r ActorRequest, key string) (model.AuditCertificate, error) {
	var out model.AuditCertificate
	if key != "" {
		if ok, err := s.Ledger.GetIdempotent("close:"+batchID, key, &out); ok || err != nil {
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
		var snap model.ReleaseSnapshot
		for _, x := range st.Snapshots {
			if x.BatchID == batchID {
				snap = x
				break
			}
		}
		if e = model.EnsureCanClose(b, &snap); e != nil {
			return e
		}
		certificateID := model.HashJSON(struct{ B, S string }{b.BatchID, snap.SnapshotID})[:24]
		s.appendEvent(st, batchID, model.NewEvent("batch_closed", r.Actor, map[string]any{"certificate_id": certificateID}))
		events := filterEvents(st.Events, batchID)
		cert, e := audit.Issue(b, snap, events, r.Actor)
		if e != nil {
			return model.Wrap("audit_issue", e)
		}
		if _, exists := st.Certificates[cert.CertificateID]; exists {
			return model.Conflict("证书已签发")
		}
		st.Certificates[cert.CertificateID] = cert
		b.Status = model.StatusClosed
		b.Version++
		b.UpdatedAt = model.Now()
		st.Batches[batchID] = b
		out = cert
		if key != "" {
			data, marshalErr := json.Marshal(out)
			if marshalErr != nil {
				return marshalErr
			}
			st.Idempotency[store.IdempotencyKey("close:"+batchID, key)] = data
		}
		return nil
	})
	return out, err
}
func filterEvents(all []model.AuditEvent, id string) []model.AuditEvent {
	out := []model.AuditEvent{}
	for _, e := range all {
		if e.BatchID == id {
			out = append(out, e)
		}
	}
	return out
}
