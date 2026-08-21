package audit

import (
	"encoding/json"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"sort"
)

func CanonicalEvent(e model.AuditEvent) []byte {
	keys := make([]string, 0, len(e.Payload))
	for k := range e.Payload {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	p := make(map[string]any, len(keys))
	for _, k := range keys {
		p[k] = e.Payload[k]
	}
	b, _ := json.Marshal(struct {
		Sequence     int64          `json:"sequence"`
		BatchID      string         `json:"batch_id"`
		Action       string         `json:"action"`
		Actor        string         `json:"actor"`
		Payload      map[string]any `json:"payload"`
		PreviousHash string         `json:"previous_hash"`
		OccurredAt   string         `json:"occurred_at"`
	}{e.Sequence, e.BatchID, e.Action, e.Actor, p, e.PreviousHash, e.OccurredAt.UTC().Format("2006-01-02T15:04:05.000Z")})
	return b
}
