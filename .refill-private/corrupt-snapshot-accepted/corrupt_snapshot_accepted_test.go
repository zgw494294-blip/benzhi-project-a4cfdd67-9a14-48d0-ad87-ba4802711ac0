package corruptsnapshotaccepted

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/store"
)

func TestLedgerRejectsSemanticallyCorruptSnapshotOnLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ledger.json")
	contents := []byte(`{
  "schema_version": 1,
  "batches": {
    "batch-corrupt": {
      "batch_id": "batch-corrupt",
      "title": "",
      "steward": "管理员",
      "policy_version": "伦理-1",
      "release_scope": ["restricted"],
      "status": "draft",
      "version": 1,
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    }
  },
  "sources": {},
  "reviews": {},
  "snapshots": {},
  "certificates": {},
  "events": [],
  "idempotency": {}
}`)
	if err := os.WriteFile(path, contents, 0600); err != nil {
		t.Fatalf("write corrupt ledger: %v", err)
	}

	if _, err := store.New(path); err == nil {
		t.Fatal("expected semantically corrupt ledger snapshot to be rejected on load")
	}
}
