package source_evidence_alias

import (
	"testing"

	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/store"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/workflow"
)

func TestAddSourceReturnDoesNotRewriteStoredEvidence(t *testing.T) {
	ledger, err := store.New("")
	if err != nil {
		t.Fatal(err)
	}
	service := workflow.New(ledger)

	batch, err := service.CreateBatch(workflow.CreateBatchRequest{
		Title:         "来源凭证隔离测试",
		Steward:       "管理员",
		PolicyVersion: "policy-1",
		ReleaseScope:  []string{"restricted"},
	}, "")
	if err != nil {
		t.Fatal(err)
	}
	source, err := service.AddSource(batch.BatchID, workflow.AddSourceRequest{
		Origin:          "研究中心",
		Description:     "匿名样本",
		ConsentRef:      "consent-original",
		Sensitivity:     "restricted",
		ContentChecksum: "checksum-original",
	}, batch.Version, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(source.EvidenceRefs) != 1 {
		t.Fatalf("expected one evidence ref, got %v", source.EvidenceRefs)
	}

	// 返回的聚合属于调用方，不得与账本存储共享底层切片。
	source.EvidenceRefs[0] = "consent-tampered"

	view, err := service.GetBatch(batch.BatchID)
	if err != nil {
		t.Fatal(err)
	}
	if len(view.Sources) != 1 {
		t.Fatalf("expected one source, got %d", len(view.Sources))
	}
	if got := view.Sources[0].EvidenceRefs; len(got) != 1 || got[0] != "consent-original" {
		t.Fatalf("stored source evidence was rewritten: %v", got)
	}
	if view.Sources[0].State != model.SourcePending {
		t.Fatalf("unexpected source state: %s", view.Sources[0].State)
	}
}
