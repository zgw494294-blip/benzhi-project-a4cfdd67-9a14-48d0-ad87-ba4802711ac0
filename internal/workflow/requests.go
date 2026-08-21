package workflow

import "github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"

type CreateBatchRequest struct {
	Title         string   `json:"title"`
	Steward       string   `json:"steward"`
	PolicyVersion string   `json:"policy_version"`
	ReleaseScope  []string `json:"release_scope"`
}
type AddSourceRequest struct {
	Origin          string `json:"origin"`
	Description     string `json:"description"`
	ConsentRef      string `json:"consent_ref"`
	Sensitivity     string `json:"sensitivity"`
	ContentChecksum string `json:"content_checksum"`
}
type AddSourcesRequest struct {
	Sources         []AddSourceRequest `json:"sources"`
	ExpectedVersion int64              `json:"expected_version"`
}
type ReviewRequest struct {
	Reviewer        string         `json:"reviewer"`
	Decision        model.Decision `json:"decision"`
	Issues          []string       `json:"issues"`
	EvidenceRefs    []string       `json:"evidence_refs"`
	ExpectedVersion int64          `json:"expected_version"`
}
type ResubmitRequest struct {
	EvidenceRefs    map[string][]string `json:"evidence_refs"`
	ExpectedVersion int64               `json:"expected_version"`
}
type ActorRequest struct {
	Actor           string `json:"actor"`
	ExpectedVersion int64  `json:"expected_version"`
}
