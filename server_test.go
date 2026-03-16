package routes

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubRouter struct {
	metrics []Metrics
	err     error
}

func (s stubRouter) Fetch(_ context.Context, _ string, _ []string) ([]Metrics, error) {
	return s.metrics, s.err
}

func TestRouteHandlerReturnsSortedResponse(t *testing.T) {
	prev := router
	router = stubRouter{
		metrics: []Metrics{
			{Duration: 500, Distance: 1800},
			{Duration: 400, Distance: 3000},
		},
	}
	t.Cleanup(func() { router = prev })

	req := httptest.NewRequest(http.MethodGet, "/route?src=13.388860,52.517037&dst=13.397634,52.529407&dst=13.428555,52.523219", nil)
	rec := httptest.NewRecorder()

	routeHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	want := "{\"source\":\"13.388860,52.517037\",\"routes\":[{\"destination\":\"13.428555,52.523219\",\"duration\":400,\"distance\":3000},{\"destination\":\"13.397634,52.529407\",\"duration\":500,\"distance\":1800}]}\n"
	if rec.Body.String() != want {
		t.Fatalf("body = %q, want %q", rec.Body.String(), want)
	}
}

func TestRouteHandlerReturnsGatewayTimeoutOnDeadlineExceeded(t *testing.T) {
	prev := router
	router = stubRouter{err: errUpstreamTimeout}
	t.Cleanup(func() { router = prev })

	req := httptest.NewRequest(http.MethodGet, "/route?src=13.388860,52.517037&dst=13.397634,52.529407", nil)
	rec := httptest.NewRecorder()

	routeHandler(rec, req)

	if rec.Code != http.StatusGatewayTimeout {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusGatewayTimeout)
	}

	if rec.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Fatalf("content-type = %q, want text/plain; charset=utf-8", rec.Header().Get("Content-Type"))
	}

	want := "osrm: upstream timeout\n"
	if rec.Body.String() != want {
		t.Fatalf("body = %q, want %q", rec.Body.String(), want)
	}
}

func TestRouteHandlerReturnsBadGatewayOnUpstreamError(t *testing.T) {
	prev := router
	router = stubRouter{err: errors.New("osrm: invalid response")}
	t.Cleanup(func() { router = prev })

	req := httptest.NewRequest(http.MethodGet, "/route?src=13.388860,52.517037&dst=13.397634,52.529407", nil)
	rec := httptest.NewRecorder()

	routeHandler(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadGateway)
	}

	want := "osrm: invalid response\n"
	if rec.Body.String() != want {
		t.Fatalf("body = %q, want %q", rec.Body.String(), want)
	}
}
