package idempotency_concurrent

import (
	"runtime"
	"sync"
	"testing"

	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/store"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/workflow"
)

func TestConcurrentCreateWithSameIdempotencyKeyCommitsOnce(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ledger, err := store.New("")
	if err != nil {
		t.Fatal(err)
	}
	service := workflow.New(ledger)
	const callers = 64
	start := make(chan struct{})
	results := make(chan model.DatasetBatch, callers)
	errs := make(chan error, callers)
	var wg sync.WaitGroup
	wg.Add(callers)
	for i := 0; i < callers; i++ {
		go func() {
			defer wg.Done()
			<-start
			batch, callErr := service.CreateBatch(workflow.CreateBatchRequest{
				Title:         "并发幂等批次",
				Steward:       "管理员",
				PolicyVersion: "伦理-1",
				ReleaseScope:  []string{"机构内部"},
			}, "same-create-key")
			results <- batch
			errs <- callErr
		}()
	}
	close(start)
	wg.Wait()
	close(results)
	close(errs)

	var firstID string
	for batch := range results {
		if firstID == "" {
			firstID = batch.BatchID
		}
		if batch.BatchID != firstID {
			t.Fatalf("幂等请求返回了不同批次: first=%q got=%q", firstID, batch.BatchID)
		}
	}
	for callErr := range errs {
		if callErr != nil {
			t.Fatalf("并发创建返回错误: %v", callErr)
		}
	}
	batches, _, _ := ledger.Count()
	if batches != 1 {
		t.Fatalf("相同 Idempotency-Key 创建了 %d 个批次", batches)
	}
}
