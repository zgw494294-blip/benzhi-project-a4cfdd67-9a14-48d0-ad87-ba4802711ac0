package snapshotslicealias

import (
	"testing"

	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/store"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/workflow"
)

func TestCreateBatchOwnsReleaseScopeSlice(t *testing.T) {
	ledger, err := store.New("")
	if err != nil {
		t.Fatalf("创建账本失败: %v", err)
	}
	service := workflow.New(ledger)
	scope := []string{"restricted", "public"}
	batch, err := service.CreateBatch(workflow.CreateBatchRequest{
		Title:         "切片所有权测试",
		Steward:       "管理员",
		PolicyVersion: "伦理-1",
		ReleaseScope:  scope,
	}, "")
	if err != nil {
		t.Fatalf("创建批次失败: %v", err)
	}

	// 请求完成后，调用方复用自己的切片；账本中的快照不应随之改变。
	scope[0] = "sensitive"
	view, err := service.GetBatch(batch.BatchID)
	if err != nil {
		t.Fatalf("读取批次失败: %v", err)
	}
	if got := view.Batch.ReleaseScope[0]; got != "restricted" {
		t.Fatalf("release_scope 被外部切片污染: got %q, want %q", got, "restricted")
	}
}
