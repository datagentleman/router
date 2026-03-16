package routes

import (
	"errors"
	"net/http/httptest"
	"testing"
)

func TestParseRequestTrimsValues(t *testing.T) {
	req := httptest.NewRequest("GET", "/route?src=%2013.388860,52.517037%20&dst=%2013.397634,52.529407%20&dst=%2013.428555,52.523219%20", nil)
	got := parseRequest(req)

	if got.Src != "13.388860,52.517037" {
		t.Fatalf("unexpected src: %q", got.Src)
	}

	if len(got.Dst) != 2 {
		t.Fatalf("unexpected dst count: %d", len(got.Dst))
	}

	if got.Dst[0] != "13.397634,52.529407" {
		t.Fatalf("unexpected first dst: %q", got.Dst[0])
	}

	if got.Dst[1] != "13.428555,52.523219" {
		t.Fatalf("unexpected second dst: %q", got.Dst[1])
	}
}

func TestRequestValidateSuccess(t *testing.T) {
	req := Request{
		Src: "13.388860,52.517037",
		Dst: []string{"13.397634,52.529407", "13.428555,52.523219"},
	}

	if err := req.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestRequestValidateFailure(t *testing.T) {
	req := Request{
		Src: "13.388860,52.517037",
		Dst: []string{"invalid"},
	}

	err := req.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want non-nil")
	}

	if !errors.Is(err, errInvalidRequest) {
		t.Fatalf("Validate() error = %v, want invalid request", err)
	}

	if !errors.Is(err, errInvalidFormat) {
		t.Fatalf("Validate() error = %v, want invalid format", err)
	}
}
