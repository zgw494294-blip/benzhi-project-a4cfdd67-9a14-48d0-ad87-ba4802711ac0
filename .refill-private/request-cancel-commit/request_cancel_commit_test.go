package requestcancelcommit

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/httpapi"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/store"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/workflow"
)

func TestCanceledCreateRequestDoesNotCommitBatch(t *testing.T) {
	ledger, err := store.New("")
	if err != nil {
		t.Fatal(err)
	}
	server := httpapi.New(workflow.New(ledger))
	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodPost, "/v1/batches", bytes.NewBufferString(`{"title":"取消测试","steward":"管理员","policy_version":"伦理-1","release_scope":["机构内部"]}`)).WithContext(ctx)
	cancel()
	server.ServeHTTP(httptest.NewRecorder(), req)

	batches, _, _ := ledger.Count()
	if batches != 0 {
		t.Fatalf("已取消的创建请求仍提交了 %d 个批次", batches)
	}
}
