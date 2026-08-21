package audit

import "github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"

func VerifyCertificate(cert model.AuditCertificate, snapshot model.ReleaseSnapshot, events []model.AuditEvent) error {
	if cert.BatchID == "" || snapshot.BatchID != cert.BatchID || cert.SnapshotID != snapshot.SnapshotID || cert.EventHash == "" || cert.AuditEntries < 1 || cert.AuditEntries != len(events) {
		return model.Invalid("审计证书字段不完整")
	}
	if err := Verify(events); err != nil {
		return model.Wrap("audit_invalid", err)
	}
	if len(events) == 0 || events[len(events)-1].EventHash != cert.EventHash {
		return model.Invalid("审计证书尾哈希不匹配")
	}
	return nil
}
