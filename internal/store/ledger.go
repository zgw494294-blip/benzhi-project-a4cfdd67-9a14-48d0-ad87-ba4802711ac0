package store

import (
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
	l.mu.Lock()
	defer l.mu.Unlock()
	next := clone(l.state)
	if err := fn(&next); err != nil {
		return err
	}
	if err := validate(next); err != nil {
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
		n.Batches[k] = cloneBatch(v)
	}
	for k, v := range s.Sources {
		n.Sources[k] = cloneSource(v)
	}
	for k, v := range s.Reviews {
		n.Reviews[k] = cloneReview(v)
	}
	for k, v := range s.Snapshots {
		n.Snapshots[k] = cloneSnapshot(v)
	}
	for k, v := range s.Certificates {
		n.Certificates[k] = v
	}
	n.Events = make([]model.AuditEvent, len(s.Events))
	for i, e := range s.Events {
		n.Events[i] = cloneEvent(e)
	}
	for k, v := range s.Idempotency {
		n.Idempotency[k] = append([]byte(nil), v...)
	}
	return n
}
func cloneBatch(b model.DatasetBatch) model.DatasetBatch {
	if b.ReleaseScope != nil {
		b.ReleaseScope = append([]string(nil), b.ReleaseScope...)
	}
	return b
}
func cloneSource(s model.SourceRecord) model.SourceRecord {
	if s.EvidenceRefs != nil {
		s.EvidenceRefs = append([]string(nil), s.EvidenceRefs...)
	}
	return s
}
func cloneReview(r model.ReviewTask) model.ReviewTask {
	if r.Issues != nil {
		r.Issues = append([]string(nil), r.Issues...)
	}
	if r.EvidenceRefs != nil {
		r.EvidenceRefs = append([]string(nil), r.EvidenceRefs...)
	}
	return r
}
func cloneSnapshot(s model.ReleaseSnapshot) model.ReleaseSnapshot {
	if s.ApprovedScope != nil {
		s.ApprovedScope = append([]string(nil), s.ApprovedScope...)
	}
	return s
}
func cloneEvent(e model.AuditEvent) model.AuditEvent {
	if e.Payload != nil {
		e.Payload = clonePayload(e.Payload)
	}
	return e
}
func clonePayload(p map[string]any) map[string]any {
	out := make(map[string]any, len(p))
	for k, v := range p {
		out[k] = cloneValue(v)
	}
	return out
}
func cloneValue(v any) any {
	switch x := v.(type) {
	case []string:
		return append([]string(nil), x...)
	case []any:
		out := make([]any, len(x))
		for i, e := range x {
			out[i] = cloneValue(e)
		}
		return out
	case map[string]any:
		return clonePayload(x)
	case []byte:
		return append([]byte(nil), x...)
	default:
		return v
	}
}
func validate(s Snapshot) error {
	if s.SchemaVersion != SchemaVersion {
		return model.Invalid("账本 schemaVersion 不受支持")
	}
	return nil
}
