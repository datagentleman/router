package routes

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestOSRMRouterFetchSuccess(t *testing.T) {
	if os.Getenv("INTEGRATION") != "1" {
		t.Skip("set INTEGRATION=1 to run live OSRM integration tests")
	}

	router := NewOSRMRouter()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	metrics, err := router.Fetch(ctx, "13.388860,52.517037", []string{
		"13.397634,52.529407",
		"13.428555,52.523219",
	})

	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}

	if len(metrics) != 2 {
		t.Fatalf("unexpected metrics count: %d", len(metrics))
	}

	for i, metric := range metrics {
		if metric.Distance <= 0 {
			t.Fatalf("metric %d distance = %f, want > 0", i, metric.Distance)
		}

		if metric.Duration <= 0 {
			t.Fatalf("metric %d duration = %f, want > 0", i, metric.Duration)
		}
	}
}

func TestOSRMRouterFetchFailure(t *testing.T) {
	if os.Getenv("INTEGRATION") != "1" {
		t.Skip("set INTEGRATION=1 to run live OSRM integration tests")
	}

	router := NewOSRMRouter()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := router.Fetch(ctx, "bad-src", []string{"bad-dst"})
	if err == nil {
		t.Fatal("Fetch() error = nil, want non-nil")
	}
}
