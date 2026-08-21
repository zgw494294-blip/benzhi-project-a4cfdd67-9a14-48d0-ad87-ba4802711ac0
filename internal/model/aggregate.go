package model

import (
	"sort"
	"strings"
)

func NormalizeScope(scope []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(scope))
	for _, v := range scope {
		v = strings.TrimSpace(v)
		if v != "" && !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	sort.Strings(out)
	return out
}
func ValidateBatch(b DatasetBatch) error {
	if b.BatchID == "" {
		return Invalid("batch_id 不能为空")
	}
	if err := ValidateBatchInput(b.Title, b.Steward, b.PolicyVersion, b.ReleaseScope); err != nil {
		return err
	}
	if b.Version < 1 {
		return Invalid("version 必须为正数")
	}
	switch b.Status {
	case StatusDraft, StatusReviewing, StatusCorrection, StatusFrozen, StatusClosed:
	default:
		return Invalid("批次状态不合法")
	}
	if !ValidTime(b.CreatedAt) || !ValidTime(b.UpdatedAt) {
		return Invalid("时间必须为 UTC")
	}
	return nil
}
func ValidateSource(s SourceRecord) error {
	if s.SourceID == "" || s.BatchID == "" {
		return Invalid("source_id 和 batch_id 不能为空")
	}
	if err := ValidateSourceInput(s.Origin, s.Description, s.ConsentRef, s.Sensitivity, s.ContentChecksum); err != nil {
		return err
	}
	switch s.State {
	case SourcePending, SourceApproved, SourceRejected, SourceNeedsCorrection:
	default:
		return Invalid("来源状态不合法")
	}
	if !ValidTime(s.CreatedAt) {
		return Invalid("来源时间必须为 UTC")
	}
	for _, ref := range s.EvidenceRefs {
		if strings.TrimSpace(ref) == "" {
			return Invalid("evidence_refs 不能包含空项")
		}
	}
	return nil
}
func ValidateReview(r ReviewTask) error {
	if r.ReviewID == "" || r.BatchID == "" || r.SourceID == "" {
		return Invalid("审核标识不能为空")
	}
	if err := ValidateReviewInput(r.Reviewer, r.Decision, r.EvidenceRefs, r.Issues); err != nil {
		return err
	}
	if r.Revision < 1 {
		return Invalid("revision 必须为正数")
	}
	if !ValidTime(r.ReviewedAt) {
		return Invalid("审核时间必须为 UTC")
	}
	return nil
}
func ValidateSnapshot(s ReleaseSnapshot) error {
	if s.SnapshotID == "" || s.BatchID == "" || s.ManifestHash == "" {
		return Invalid("发布快照字段不完整")
	}
	if s.SourceCount < 1 || len(s.ApprovedScope) == 0 {
		return Invalid("发布快照内容不完整")
	}
	if !ValidTime(s.FrozenAt) {
		return Invalid("冻结时间必须为 UTC")
	}
	return nil
}
func ValidateCertificate(c AuditCertificate) error {
	if c.CertificateID == "" || c.BatchID == "" || c.SnapshotID == "" || c.EventHash == "" {
		return Invalid("审计证书字段不完整")
	}
	if c.AuditEntries < 1 || c.IssuedBy == "" || !ValidTime(c.IssuedAt) {
		return Invalid("审计证书内容不完整")
	}
	return nil
}
func ValidateAggregate(b DatasetBatch, sources []SourceRecord, reviews []ReviewTask, snapshot *ReleaseSnapshot, cert *AuditCertificate) error {
	if err := ValidateBatch(b); err != nil {
		return err
	}
	ids := map[string]bool{}
	for _, s := range sources {
		if err := ValidateSource(s); err != nil {
			return err
		}
		if s.BatchID != b.BatchID {
			return Invalid("来源批次不一致")
		}
		if ids[s.SourceID] {
			return Invalid("来源 ID 重复")
		}
		ids[s.SourceID] = true
	}
	for _, r := range reviews {
		if err := ValidateReview(r); err != nil {
			return err
		}
		if r.BatchID != b.BatchID || !ids[r.SourceID] {
			return Invalid("审核引用不存在的来源")
		}
	}
	if snapshot != nil {
		if err := ValidateSnapshot(*snapshot); err != nil {
			return err
		}
		if snapshot.BatchID != b.BatchID {
			return Invalid("快照批次不一致")
		}
	}
	if cert != nil {
		if err := ValidateCertificate(*cert); err != nil {
			return err
		}
		if cert.BatchID != b.BatchID {
			return Invalid("证书批次不一致")
		}
	}
	if b.Status == StatusFrozen || b.Status == StatusClosed {
		if snapshot == nil {
			return Invalid("冻结或关闭批次必须有快照")
		}
	}
	if b.Status == StatusClosed && cert == nil {
		return Invalid("关闭批次必须有证书")
	}
	return nil
}
