package store

import (
	"encoding/json"
	"errors"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"os"
	"path/filepath"
)

func (l *Ledger) load() error {
	if l.path == "" {
		return nil
	}
	b, err := os.ReadFile(l.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	var s Snapshot
	if err = json.Unmarshal(b, &s); err != nil {
		return model.Err("ledger_corrupt", "账本 JSON 损坏")
	}
	if err = validate(s); err != nil {
		return err
	}
	if s.Batches == nil {
		s.Batches = map[string]model.DatasetBatch{}
	}
	if s.Sources == nil {
		s.Sources = map[string]model.SourceRecord{}
	}
	if s.Reviews == nil {
		s.Reviews = map[string]model.ReviewTask{}
	}
	if s.Snapshots == nil {
		s.Snapshots = map[string]model.ReleaseSnapshot{}
	}
	if s.Certificates == nil {
		s.Certificates = map[string]model.AuditCertificate{}
	}
	if s.Idempotency == nil {
		s.Idempotency = map[string][]byte{}
	}
	l.state = s
	return nil
}
func (l *Ledger) persist(s Snapshot) error {
	if l.path == "" {
		return nil
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(l.path)
	if err = os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "ledger-*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if _, err = tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err = tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err = tmp.Close(); err != nil {
		return err
	}
	return os.Rename(name, l.path)
}
func (l *Ledger) Path() string { return l.path }
