package audit

import (
	"fmt"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
)

type Chain struct{ events []model.AuditEvent }

func New(events []model.AuditEvent) *Chain {
	return &Chain{events: append([]model.AuditEvent(nil), events...)}
}
func (c *Chain) Append(batchID string, event model.Event) model.AuditEvent {
	var prev string
	if len(c.events) > 0 {
		prev = c.events[len(c.events)-1].EventHash
	}
	e := model.AuditEvent{Sequence: int64(len(c.events) + 1), BatchID: batchID, Action: event.Action, Actor: event.Actor, Payload: event.Payload, PreviousHash: prev, OccurredAt: model.Now()}
	e.EventHash = model.HashJSON(struct {
		Sequence               int64
		BatchID, Action, Actor string
		Payload                map[string]any
		PreviousHash           string
		OccurredAt             string
	}{e.Sequence, e.BatchID, e.Action, e.Actor, e.Payload, e.PreviousHash, e.OccurredAt.Format("2006-01-02T15:04:05.000Z")})
	c.events = append(c.events, e)
	return e
}
func (c *Chain) Events() []model.AuditEvent { return append([]model.AuditEvent(nil), c.events...) }
func Verify(events []model.AuditEvent) error {
	var prev string
	for i, e := range events {
		if e.Sequence != int64(i+1) || e.PreviousHash != prev {
			return fmt.Errorf("审计链顺序或前序哈希错误")
		}
		expected := model.HashJSON(struct {
			Sequence               int64
			BatchID, Action, Actor string
			Payload                map[string]any
			PreviousHash           string
			OccurredAt             string
		}{e.Sequence, e.BatchID, e.Action, e.Actor, e.Payload, e.PreviousHash, e.OccurredAt.Format("2006-01-02T15:04:05.000Z")})
		if expected != e.EventHash {
			return fmt.Errorf("审计事件哈希不匹配")
		}
		prev = e.EventHash
	}
	return nil
}
