package httpapi

import (
	"context"
	"errors"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
	"net/http"
)

func respondError(w http.ResponseWriter, err error) {
	if err == context.Canceled {
		writeJSON(w, 499, map[string]any{"error": map[string]string{"code": "request_cancelled", "message": "客户端已取消请求"}})
		return
	}
	if err == context.DeadlineExceeded {
		writeJSON(w, http.StatusGatewayTimeout, map[string]any{"error": map[string]string{"code": "request_timeout", "message": "请求超时"}})
		return
	}
	status := http.StatusBadRequest
	var d *model.DomainError
	if errors.As(err, &d) {
		switch d.Code {
		case "not_found":
			status = http.StatusNotFound
		case "conflict", "version_mismatch":
			status = http.StatusConflict
		case "ledger_corrupt", "audit_invalid", "audit_issue":
			status = http.StatusInternalServerError
		}
		writeJSON(w, status, map[string]any{"error": map[string]string{"code": d.Code, "message": d.Message}})
		return
	}
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": "internal_error", "message": err.Error()}})
}
