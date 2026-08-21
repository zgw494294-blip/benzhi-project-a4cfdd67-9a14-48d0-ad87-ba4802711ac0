package model

import "strings"

func ValidateBatchInput(title, steward, policy string, scope []string) error {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(steward) == "" || strings.TrimSpace(policy) == "" {
		return Invalid("title、steward、policy_version 不能为空")
	}
	if len(scope) == 0 {
		return Invalid("release_scope 至少包含一项")
	}
	for _, v := range scope {
		if strings.TrimSpace(v) == "" {
			return Invalid("release_scope 不能包含空项")
		}
	}
	return nil
}
func ValidateSourceInput(origin, desc, consent, sensitivity, checksum string) error {
	if strings.TrimSpace(origin) == "" || strings.TrimSpace(desc) == "" || strings.TrimSpace(consent) == "" || strings.TrimSpace(checksum) == "" {
		return Invalid("来源字段和凭证引用不能为空")
	}
	if sensitivity != "public" && sensitivity != "restricted" && sensitivity != "sensitive" {
		return Invalid("sensitivity 必须为 public、restricted 或 sensitive")
	}
	return nil
}
func ValidateReviewInput(reviewer string, decision Decision, evidence, issues []string) error {
	if strings.TrimSpace(reviewer) == "" {
		return Invalid("reviewer 不能为空")
	}
	if decision != DecisionApprove && decision != DecisionReject && decision != DecisionCorrection {
		return Invalid("decision 不合法")
	}
	if decision == DecisionApprove && len(issues) > 0 {
		return Invalid("通过审核不能携带 issues")
	}
	if decision != DecisionApprove && len(evidence) == 0 {
		return Invalid("未通过审核必须提供 evidence_refs")
	}
	return nil
}
func EnsureCanAddSource(b DatasetBatch) error {
	if b.Status != StatusDraft && b.Status != StatusCorrection {
		return Conflict("当前状态不可登记来源")
	}
	return nil
}
func EnsureCanReview(b DatasetBatch, sources []SourceRecord, sourceID string) error {
	if b.Status != StatusDraft && b.Status != StatusReviewing && b.Status != StatusCorrection {
		return Conflict("当前状态不可审核")
	}
	for _, s := range sources {
		if s.SourceID == sourceID {
			if s.State == SourceApproved {
				return Conflict("来源已经通过审核")
			}
			return nil
		}
	}
	return NotFound("source", sourceID)
}
func AllApproved(sources []SourceRecord) bool {
	return len(sources) > 0 && all(sources, func(s SourceRecord) bool { return s.State == SourceApproved })
}
func all[T any](items []T, pred func(T) bool) bool {
	for _, x := range items {
		if !pred(x) {
			return false
		}
	}
	return true
}
func EnsureCanFreeze(b DatasetBatch, sources []SourceRecord) error {
	if b.Status == StatusFrozen || b.Status == StatusClosed {
		return Conflict("批次已经冻结或关闭")
	}
	if !AllApproved(sources) {
		return Conflict("仍有来源未通过审核")
	}
	return nil
}
func EnsureCanClose(b DatasetBatch, snapshot *ReleaseSnapshot) error {
	if b.Status != StatusFrozen {
		return Conflict("只有冻结批次可以关闭")
	}
	if snapshot == nil {
		return Conflict("缺少发布快照")
	}
	return nil
}
