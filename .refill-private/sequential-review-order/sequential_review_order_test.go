package sequential_review_order

import (
	"testing"

	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/store"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/workflow"
)

func TestReviewRejectsOutOfOrderSource(t *testing.T) {
	ledger, err := store.New("")
	if err != nil {
		t.Fatalf("创建账本失败: %v", err)
	}
	service := workflow.New(ledger)
	batch, err := service.CreateBatch(workflow.CreateBatchRequest{
		Title:         "顺序审核测试",
		Steward:       "管理员",
		PolicyVersion: "伦理-1",
		ReleaseScope:  []string{"restricted"},
	}, "")
	if err != nil {
		t.Fatalf("创建批次失败: %v", err)
	}
	sources, err := service.AddSources(batch.BatchID, workflow.AddSourcesRequest{
		Sources: []workflow.AddSourceRequest{
			{Origin: "实验室甲", Description: "匿名样本甲", ConsentRef: "consent-a", Sensitivity: "restricted", ContentChecksum: "checksum-a"},
			{Origin: "实验室乙", Description: "匿名样本乙", ConsentRef: "consent-b", Sensitivity: "restricted", ContentChecksum: "checksum-b"},
		},
		ExpectedVersion: batch.Version,
	}, batch.Version, "")
	if err != nil {
		t.Fatalf("批量登记来源失败: %v", err)
	}
	if len(sources) != 2 {
		t.Fatalf("来源数量错误: %d", len(sources))
	}

	_, err = service.Review(batch.BatchID, sources[1].SourceID, workflow.ReviewRequest{
		Reviewer:        "审查员",
		Decision:        model.DecisionApprove,
		EvidenceRefs:    []string{"evidence-b"},
		ExpectedVersion: batch.Version + 1,
	}, "")
	if err == nil {
		t.Fatalf("TestReviewRejectsOutOfOrderSource: 第二个来源在第一个来源审核前被接受")
	}
	if !model.IsCode(err, "conflict") {
		t.Fatalf("越序审核应返回 conflict，实际为: %v", err)
	}
	view, viewErr := service.GetBatch(batch.BatchID)
	if viewErr != nil {
		t.Fatalf("读取批次失败: %v", viewErr)
	}
	if len(view.Reviews) != 0 {
		t.Fatalf("越序审核不应写入审核记录: %#v", view.Reviews)
	}
}
