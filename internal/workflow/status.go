package workflow

import (
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/audit"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"sort"
)

type StatusReport struct {
	Batch             model.DatasetBatch `json:"batch"`
	Approved          int                `json:"approved"`
	Total             int                `json:"total"`
	AuditValid        bool               `json:"audit_valid"`
	AuditActions      []string           `json:"audit_actions"`
	Pending           int                `json:"pending"`
	Correction        int                `json:"correction"`
	PendingCount      int                `json:"pending_count"`
	CorrectionCount   int                `json:"correction_count"`
	PendingSources    []string           `json:"pending_sources"`
	CorrectionSources []string           `json:"correction_sources"`
}

func (s *Service) Status(id string) (StatusReport, error) {
	v, e := s.GetBatch(id)
	if e != nil {
		return StatusReport{}, e
	}
	a, t := model.ReviewCoverage(v.Sources, v.Reviews)
	pending := []string{}
	correction := []string{}
	for _, src := range v.Sources {
		switch src.State {
		case model.SourcePending:
			pending = append(pending, src.SourceID)
		case model.SourceNeedsCorrection, model.SourceRejected:
			correction = append(correction, src.SourceID)
		}
	}
	sort.Strings(pending)
	sort.Strings(correction)
	return StatusReport{Batch: v.Batch, Approved: a, Total: t, Pending: len(pending), Correction: len(correction), PendingCount: len(pending), CorrectionCount: len(correction), PendingSources: pending, CorrectionSources: correction, AuditValid: audit.Verify(v.Audit) == nil, AuditActions: audit.EventActions(v.Audit)}, nil
}
func (s *Service) Verify(id string) error {
	v, e := s.GetBatch(id)
	if e != nil {
		return e
	}
	if v.Batch.Status != model.StatusFrozen && v.Batch.Status != model.StatusClosed {
		return model.Conflict("批次尚未封存")
	}
	if e = audit.Verify(v.Audit); e != nil {
		return model.Wrap("audit_invalid", e)
	}
	if v.Certificate == nil || v.Snapshot == nil {
		return model.Wrap("audit_invalid", model.Invalid("封存产物缺少快照或证书"))
	}
	if e = audit.VerifyCertificate(*v.Certificate, *v.Snapshot, v.Audit); e != nil {
		return model.Wrap("audit_invalid", e)
	}
	return nil
}
