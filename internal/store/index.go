package store

import (
	"encoding/json"
	"fmt"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
)

func IdempotencyKey(scope, key string) string { return fmt.Sprintf("%s:%s", scope, key) }
func (l *Ledger) GetIdempotent(scope, key string, out any) (bool, error) {
	found := false
	err := l.Read(func(s Snapshot) error {
		b, ok := s.Idempotency[IdempotencyKey(scope, key)]
		if !ok {
			return nil
		}
		found = true
		return json.Unmarshal(b, out)
	})
	return found, err
}
func (l *Ledger) PutIdempotent(scope, key string, value any) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return l.Update(func(s *Snapshot) error {
		k := IdempotencyKey(scope, key)
		if _, ok := s.Idempotency[k]; !ok {
			s.Idempotency[k] = b
		}
		return nil
	})
}
func (l *Ledger) Events(batchID string) ([]model.AuditEvent, error) {
	out := []model.AuditEvent{}
	err := l.Read(func(s Snapshot) error {
		for _, e := range s.Events {
			if e.BatchID == batchID {
				out = append(out, e)
			}
		}
		return nil
	})
	return out, err
}
