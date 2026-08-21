package audit

import (
	"fmt"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
)

func Issue(batch model.DatasetBatch, snapshot model.ReleaseSnapshot, events []model.AuditEvent, issuer string) (model.AuditCertificate, error) {
	if err := Verify(events); err != nil {
		return model.AuditCertificate{}, err
	}
	if len(events) == 0 {
		return model.AuditCertificate{}, fmt.Errorf("没有审计事件")
	}
	return model.AuditCertificate{CertificateID: model.HashJSON(struct{ B, S string }{batch.BatchID, snapshot.SnapshotID})[:24], BatchID: batch.BatchID, SnapshotID: snapshot.SnapshotID, EventHash: events[len(events)-1].EventHash, AuditEntries: len(events), IssuedAt: model.Now(), IssuedBy: issuer}, nil
}
