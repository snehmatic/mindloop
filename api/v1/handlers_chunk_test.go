package v1_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/snehmatic/mindloop/api/v1"
)

func TestHandleChunk_InvalidType(t *testing.T) {
	// Initialize a dummy handler (we can pass nil for services just to test validation)
	mlh := v1.NewMindloopHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	form := url.Values{}
	form.Add("type", "invalid")
	form.Add("id", "1")

	req := httptest.NewRequest("POST", "/api/v1/chunk", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	mlh.HandleChunk(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
