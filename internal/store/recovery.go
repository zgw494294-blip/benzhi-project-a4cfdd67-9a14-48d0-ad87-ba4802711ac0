package store

import (
	"os"
	"path/filepath"
)

func (l *Ledger) Backup() error {
	if l.path == "" {
		return nil
	}
	data, err := os.ReadFile(l.path)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Clean(l.path)+".bak", data, 0600)
}
func (l *Ledger) Exists() bool {
	if l.path == "" {
		return false
	}
	_, err := os.Stat(l.path)
	return err == nil
}
func (l *Ledger) Reload() error { l.mu.Lock(); defer l.mu.Unlock(); return l.load() }
