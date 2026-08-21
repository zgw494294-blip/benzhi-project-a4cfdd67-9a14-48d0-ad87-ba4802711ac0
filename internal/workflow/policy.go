package workflow

import (
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"sort"
	"strings"
)

type PolicyResult struct {
	Allowed bool     `json:"allowed"`
	Reasons []string `json:"reasons,omitempty"`
}

func EvaluateSourcePolicy(b model.DatasetBatch, s model.SourceRecord) PolicyResult {
	r := PolicyResult{Allowed: true}
	if !model.ScopeAllows(b.ReleaseScope, s.Sensitivity) {
		r.Allowed = false
		r.Reasons = append(r.Reasons, "敏感级别超出发布范围")
	}
	if strings.TrimSpace(s.ConsentRef) == "" {
		r.Allowed = false
		r.Reasons = append(r.Reasons, "缺少同意或授权凭证")
	}
	if strings.TrimSpace(s.ContentChecksum) == "" {
		r.Allowed = false
		r.Reasons = append(r.Reasons, "缺少内容哈希")
	}
	return r
}
func (s *Service) Policy(id string) (map[string]PolicyResult, error) {
	v, e := s.GetBatch(id)
	if e != nil {
		return nil, e
	}
	out := map[string]PolicyResult{}
	for _, src := range v.Sources {
		out[src.SourceID] = EvaluateSourcePolicy(v.Batch, src)
	}
	return out, nil
}

type PolicyCheck struct {
	SourceID string   `json:"source_id"`
	Allowed  bool     `json:"allowed"`
	Reasons  []string `json:"reasons,omitempty"`
}

type PolicyReport struct {
	BatchID      string        `json:"batch_id"`
	Results      []PolicyCheck `json:"results"`
	Allowed      int           `json:"allowed"`
	Blocked      int           `json:"blocked"`
	AllowedCount int           `json:"allowed_count"`
	BlockedCount int           `json:"blocked_count"`
}

func (s *Service) PolicyPrecheck(id string) (PolicyReport, error) {
	v, err := s.GetBatch(id)
	if err != nil {
		return PolicyReport{}, err
	}
	results := make([]PolicyCheck, 0, len(v.Sources))
	for _, src := range v.Sources {
		result := EvaluateSourcePolicy(v.Batch, src)
		results = append(results, PolicyCheck{SourceID: src.SourceID, Allowed: result.Allowed, Reasons: append([]string(nil), result.Reasons...)})
	}
	sort.Slice(results, func(i, j int) bool { return results[i].SourceID < results[j].SourceID })
	report := PolicyReport{BatchID: id, Results: results}
	for _, result := range results {
		if result.Allowed {
			report.Allowed++
			report.AllowedCount++
		} else {
			report.Blocked++
			report.BlockedCount++
		}
	}
	return report, nil
}
func MergeEvidence(old []string, newRefs []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, x := range append(old, newRefs...) {
		x = strings.TrimSpace(x)
		if x != "" && !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	return out
}
func ResolveIssues(issues []string, evidence []string) []string {
	if len(evidence) == 0 {
		return append([]string(nil), issues...)
	}
	return []string{}
}
