package workflow

import (
	"encoding/json"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/audit"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"sort"
)

type Export struct {
	Batch        model.DatasetBatch      `json:"batch"`
	Sources      []model.SourceRecord    `json:"sources"`
	Snapshot     *model.ReleaseSnapshot  `json:"snapshot,omitempty"`
	Certificate  *model.AuditCertificate `json:"certificate,omitempty"`
	AuditSummary audit.Summary           `json:"audit_summary"`
}

func (s *Service) Export(id string) ([]byte, error) {
	v, e := s.GetBatch(id)
	if e != nil {
		return nil, e
	}
	if v.Batch.Status != model.StatusFrozen && v.Batch.Status != model.StatusClosed {
		return nil, model.Conflict("批次尚未冻结，不能导出封存产物")
	}
	if v.Snapshot == nil {
		return nil, model.Wrap("audit_invalid", model.Invalid("封存快照不存在"))
	}
	if err := audit.Verify(v.Audit); err != nil {
		return nil, model.Wrap("audit_invalid", err)
	}
	if v.Batch.Status == model.StatusClosed {
		if v.Certificate == nil {
			return nil, model.Wrap("audit_invalid", model.Invalid("关闭批次缺少审计证书"))
		}
		if err := audit.VerifyCertificate(*v.Certificate, *v.Snapshot, v.Audit); err != nil {
			return nil, model.Wrap("audit_invalid", err)
		}
	}
	sort.Slice(v.Sources, func(i, j int) bool { return v.Sources[i].SourceID < v.Sources[j].SourceID })
	if model.ManifestHash(v.Batch, v.Sources) != v.Snapshot.ManifestHash {
		return nil, model.Wrap("audit_invalid", model.Invalid("manifest_hash 不匹配"))
	}
	x := Export{Batch: v.Batch, Sources: v.Sources, Snapshot: v.Snapshot, Certificate: v.Certificate, AuditSummary: audit.Summarize(id, v.Audit)}
	return json.MarshalIndent(x, "", "  ")
}
func (s *Service) ApprovedManifest(id string) (string, error) {
	v, e := s.GetBatch(id)
	if e != nil {
		return "", e
	}
	if v.Snapshot == nil {
		return "", model.Conflict("批次尚未冻结")
	}
	return v.Snapshot.ManifestHash, nil
}
func (s *Service) Pending(id string) ([]string, error) {
	v, e := s.GetBatch(id)
	if e != nil {
		return nil, e
	}
	return model.PendingSources(v.Sources), nil
}
