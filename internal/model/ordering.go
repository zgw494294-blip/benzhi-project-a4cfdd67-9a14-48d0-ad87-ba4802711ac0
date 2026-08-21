package model

import "sort"

func SortSources(items []SourceRecord) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].CreatedAt.Equal(items[j].CreatedAt) {
			return items[i].SourceID < items[j].SourceID
		}
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})
}
func SortReviews(items []ReviewTask) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].ReviewedAt.Equal(items[j].ReviewedAt) {
			return items[i].ReviewID < items[j].ReviewID
		}
		return items[i].ReviewedAt.Before(items[j].ReviewedAt)
	})
}
func LatestReview(reviews []ReviewTask, sourceID string) (ReviewTask, bool) {
	var best ReviewTask
	ok := false
	for _, r := range reviews {
		if r.SourceID == sourceID && (!ok || r.Revision > best.Revision) {
			best = r
			ok = true
		}
	}
	return best, ok
}
func ReviewCoverage(sources []SourceRecord, reviews []ReviewTask) (int, int) {
	approved := 0
	for _, s := range sources {
		r, ok := LatestReview(reviews, s.SourceID)
		if ok && r.Decision == DecisionApprove {
			approved++
		}
	}
	return approved, len(sources)
}

func SourceIDs(items []SourceRecord) []string {
	out := make([]string, 0, len(items))
	for _, s := range items {
		out = append(out, s.SourceID)
	}
	sort.Strings(out)
	return out
}
