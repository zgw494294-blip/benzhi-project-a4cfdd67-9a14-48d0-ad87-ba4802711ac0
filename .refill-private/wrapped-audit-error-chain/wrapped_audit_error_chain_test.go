package wrapped_audit_error_chain_test

import (
	"errors"
	"testing"

	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/store"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/workflow"
)

func TestVerifyPreservesWrappedAuditErrorCause(t *testing.T) {
	ledger, err := store.New("")
	if err != nil {
		t.Fatalf("create ledger: %v", err)
	}
	service := workflow.New(ledger)
	batch, err := service.CreateBatch(workflow.CreateBatchRequest{
		Title:         "错误链测试",
		Steward:       "管理员",
		PolicyVersion: "policy-v1",
		ReleaseScope:  []string{"机构内部"},
	}, "")
	if err != nil {
		t.Fatalf("create batch: %v", err)
	}
	source, err := service.AddSource(batch.BatchID, workflow.AddSourceRequest{
		Origin:          "实验室",
		Description:     "匿名样本",
		ConsentRef:      "consent-1",
		Sensitivity:     "restricted",
		ContentChecksum: "checksum-1",
	}, batch.Version, "")
	if err != nil {
		t.Fatalf("add source: %v", err)
	}
	view, err := service.GetBatch(batch.BatchID)
	if err != nil {
		t.Fatalf("read batch: %v", err)
	}
	if _, err = service.Review(batch.BatchID, source.SourceID, workflow.ReviewRequest{
		Reviewer:        "审查员",
		Decision:        model.DecisionApprove,
		EvidenceRefs:    []string{"evidence-1"},
		ExpectedVersion: view.Batch.Version,
	}, ""); err != nil {
		t.Fatalf("review source: %v", err)
	}
	view, err = service.GetBatch(batch.BatchID)
	if err != nil {
		t.Fatalf("read reviewed batch: %v", err)
	}
	if _, err = service.Freeze(batch.BatchID, workflow.ActorRequest{Actor: "发布负责人", ExpectedVersion: view.Batch.Version}, ""); err != nil {
		t.Fatalf("freeze batch: %v", err)
	}
	view, err = service.GetBatch(batch.BatchID)
	if err != nil {
		t.Fatalf("read frozen batch: %v", err)
	}
	if _, err = service.Close(batch.BatchID, workflow.ActorRequest{Actor: "发布负责人", ExpectedVersion: view.Batch.Version}, ""); err != nil {
		t.Fatalf("close batch: %v", err)
	}

	if err = ledger.Update(func(snapshot *store.Snapshot) error {
		for id, cert := range snapshot.Certificates {
			cert.EventHash = ""
			snapshot.Certificates[id] = cert
		}
		return nil
	}); err != nil {
		t.Fatalf("corrupt certificate for verification scenario: %v", err)
	}

	err = service.Verify(batch.BatchID)
	var outer *model.DomainError
	if !errors.As(err, &outer) || outer.Code != "audit_invalid" {
		t.Fatalf("Verify returned %T %v, want audit_invalid wrapper", err, err)
	}
	inner := errors.Unwrap(err)
	var detail *model.DomainError
	if !errors.As(inner, &detail) || detail.Code != "invalid_request" {
		t.Fatalf("Verify lost wrapped cause: got %T %v, want invalid_request cause", inner, inner)
	}
}
