package audit_payload_alias

import (
	"testing"

	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/store"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/workflow"
)

func TestAuditReadOwnsPayloadMap(t *testing.T) {
	ledger, err := store.New("")
	if err != nil {
		t.Fatal(err)
	}
	service := workflow.New(ledger)
	batch, err := service.CreateBatch(workflow.CreateBatchRequest{
		Title:         "开放数据集",
		Steward:       "管理员",
		PolicyVersion: "policy-1",
		ReleaseScope:  []string{"public"},
	}, "")
	if err != nil {
		t.Fatal(err)
	}

	first, err := service.GetAudit(batch.BatchID)
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 1 {
		t.Fatalf("expected one audit event, got %d", len(first))
	}
	first[0].Payload["title"] = "被篡改的标题"

	second, err := service.GetAudit(batch.BatchID)
	if err != nil {
		t.Fatal(err)
	}
	if got := second[0].Payload["title"]; got != "开放数据集" {
		t.Fatalf("audit payload was mutated through a returned map: got %v", got)
	}
	if second[0].Action != "batch_created" || second[0].Actor != "管理员" {
		t.Fatalf("unexpected audit event: %+v", second[0])
	}
}
