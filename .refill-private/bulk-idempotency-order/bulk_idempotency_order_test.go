package bulk_idempotency_order

import (
	"fmt"
	"testing"

	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/store"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/workflow"
)

func TestBulkSourceIdempotencyReplaysCanonicalOrder(t *testing.T) {
	ledger, err := store.New("")
	if err != nil {
		t.Fatal(err)
	}
	service := workflow.New(ledger)
	batch, err := service.CreateBatch(workflow.CreateBatchRequest{
		Title:         "批量顺序",
		Steward:       "管理员",
		PolicyVersion: "伦理-1",
		ReleaseScope:  []string{"机构内部"},
	}, "")
	if err != nil {
		t.Fatal(err)
	}
	request := workflow.AddSourcesRequest{ExpectedVersion: batch.Version}
	for i := 0; i < 16; i++ {
		request.Sources = append(request.Sources, workflow.AddSourceRequest{
			Origin:          fmt.Sprintf("来源-%02d", i),
			Description:     fmt.Sprintf("描述-%02d", i),
			ConsentRef:      fmt.Sprintf("consent-%02d", i),
			Sensitivity:     "restricted",
			ContentChecksum: fmt.Sprintf("checksum-%02d", i),
		})
	}
	first, err := service.AddSources(batch.BatchID, request, batch.Version, "bulk-order-key")
	if err != nil {
		t.Fatal(err)
	}
	replayed, err := service.AddSources(batch.BatchID, request, batch.Version, "bulk-order-key")
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != len(replayed) {
		t.Fatalf("replayed count = %d, want %d", len(replayed), len(first))
	}
	for i := range first {
		if first[i].SourceID != replayed[i].SourceID {
			t.Fatalf("idempotent replay order differs at index %d: first=%q replayed=%q", i, first[i].SourceID, replayed[i].SourceID)
		}
	}
}
