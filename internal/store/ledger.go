package store

import (
	"context"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"sync"
)

const SchemaVersion = 1

type Snapshot struct {
	SchemaVersion int                               `json:"schema_version"`
	Batches       map[string]model.DatasetBatch     `json:"batches"`
	Sources       map[string]model.SourceRecord     `json:"sources"`
	Reviews       map[string]model.ReviewTask       `json:"reviews"`
	Snapshots     map[string]model.ReleaseSnapshot  `json:"snapshots"`
	Certificates  map[string]model.AuditCertificate `json:"certificates"`
	Events        []model.AuditEvent                `json:"events"`
	Idempotency   map[string][]byte                 `json:"idempotency"`
}
type Ledger struct {
	mu    sync.RWMutex
	path  string
	state Snapshot
}

func EmptySnapshot() Snapshot {
	return Snapshot{SchemaVersion: SchemaVersion, Batches: map[string]model.DatasetBatch{}, Sources: map[string]model.SourceRecord{}, Reviews: map[string]model.ReviewTask{}, Snapshots: map[string]model.ReleaseSnapshot{}, Certificates: map[string]model.AuditCertificate{}, Events: []model.AuditEvent{}, Idempotency: map[string][]byte{}}
}
func New(path string) (*Ledger, error) {
	l := &Ledger{path: path, state: EmptySnapshot()}
	if err := l.load(); err != nil {
		return nil, err
	}
	return l, nil
}
func (l *Ledger) Read(fn func(Snapshot) error) error {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return fn(clone(l.state))
}
func (l *Ledger) Update(fn func(*Snapshot) error) error {
	return l.UpdateContext(context.Background(), fn)
}
func (l *Ledger) UpdateContext(ctx context.Context, fn func(*Snapshot) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	next := clone(l.state)
	if err := fn(&next); err != nil {
		return err
	}
	if err := validate(next); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := l.persist(next); err != nil {
		return err
	}
	l.state = next
	return nil
}
func clone(s Snapshot) Snapshot {
	n := EmptySnapshot()
	n.SchemaVersion = s.SchemaVersion
	for k, v := range s.Batches {
		n.Batches[k] = v
	}
	for k, v := range s.Sources {
		n.Sources[k] = v
	}
	for k, v := range s.Reviews {
		n.Reviews[k] = v
	}
	for k, v := range s.Snapshots {
		n.Snapshots[k] = v
	}
	for k, v := range s.Certificates {
		n.Certificates[k] = v
	}
	n.Events = append([]model.AuditEvent(nil), s.Events...)
	for k, v := range s.Idempotency {
		n.Idempotency[k] = append([]byte(nil), v...)
	}
	return n
}
func validate(s Snapshot) error {
	if s.SchemaVersion != SchemaVersion {
		return model.Invalid("账本 schemaVersion 不受支持")
	}
	return nil
}
