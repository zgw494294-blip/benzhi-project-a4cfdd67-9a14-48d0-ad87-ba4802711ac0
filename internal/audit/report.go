package audit

import (
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"sort"
)

type Summary struct {
	BatchID     string   `json:"batch_id"`
	EventCount  int      `json:"event_count"`
	FirstAction string   `json:"first_action"`
	LastAction  string   `json:"last_action"`
	Actors      []string `json:"actors"`
	Valid       bool     `json:"valid"`
}

func Summarize(batchID string, events []model.AuditEvent) Summary {
	s := Summary{BatchID: batchID, EventCount: len(events)}
	seen := map[string]bool{}
	for _, e := range events {
		if s.FirstAction == "" {
			s.FirstAction = e.Action
		}
		s.LastAction = e.Action
		if !seen[e.Actor] {
			seen[e.Actor] = true
			s.Actors = append(s.Actors, e.Actor)
		}
	}
	sort.Strings(s.Actors)
	s.Valid = Verify(events) == nil
	return s
}
func EventActions(events []model.AuditEvent) []string {
	out := make([]string, 0, len(events))
	for _, e := range events {
		out = append(out, e.Action)
	}
	return out
}
func TailHash(events []model.AuditEvent) string {
	if len(events) == 0 {
		return ""
	}
	return events[len(events)-1].EventHash
}
