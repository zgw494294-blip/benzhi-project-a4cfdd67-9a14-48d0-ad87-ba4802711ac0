package reviewissuesalias

import (
	"testing"

	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/store"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/workflow"
)

func TestReviewRequestIssuesDoNotRewriteAudit(t *testing.T) {
	ledger, err := store.New("")
	if err != nil {
		t.Fatal(err)
	}
	service := workflow.New(ledger)

	batch, err := service.CreateBatch(workflow.CreateBatchRequest{
		Title:         "受控数据集",
		Steward:       "管理员",
		PolicyVersion: "伦理-1",
		ReleaseScope:  []string{"机构内部"},
	}, "")
	if err != nil {
		t.Fatal(err)
	}
	source, err := service.AddSource(batch.BatchID, workflow.AddSourceRequest{
		Origin:          "实验室",
		Description:     "匿名样本",
		ConsentRef:      "consent-1",
		Sensitivity:     "restricted",
		ContentChecksum: "checksum-1",
	}, batch.Version, "")
	if err != nil {
		t.Fatal(err)
	}

	issues := []string{"缺少伦理批件"}
	_, err = service.Review(batch.BatchID, source.SourceID, workflow.ReviewRequest{
		Reviewer:        "审查员",
		Decision:        model.DecisionCorrection,
		Issues:          issues,
		EvidenceRefs:    []string{"evidence-1"},
		ExpectedVersion: batch.Version + 1,
	}, "")
	if err != nil {
		t.Fatal(err)
	}

	issues[0] = "调用方复用后的内容"
	auditEvents, err := service.GetAudit(batch.BatchID)
	if err != nil {
		t.Fatal(err)
	}
	for _, event := range auditEvents {
		if event.Action != "source_reviewed" {
			continue
		}
		payloadIssues, ok := event.Payload["issues"].([]string)
		if !ok || len(payloadIssues) != 1 {
			t.Fatalf("审计事件缺少 issues: %#v", event.Payload["issues"])
		}
		if payloadIssues[0] != "缺少伦理批件" {
			t.Fatalf("审计事件被请求切片改写为 %q", payloadIssues[0])
		}
		return
	}
	t.Fatal("未找到 source_reviewed 审计事件")
}
