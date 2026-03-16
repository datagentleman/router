package routes

import "testing"

func TestNewResponseSortsByDurationThenDistance(t *testing.T) {
	src := "13.388860,52.517037"

	dst := []string{
		"13.397634,52.529407",
		"13.428555,52.523219",
		"13.418555,52.521219",
	}

	metrics := []Metrics{
		{Duration: 500, Distance: 1800},
		{Duration: 400, Distance: 3000},
		{Duration: 400, Distance: 2000},
	}

	resp := newResponse(src, dst, metrics)

	if resp.Source != src {
		t.Fatalf("unexpected source: %q", resp.Source)
	}

	if len(resp.Routes) != 3 {
		t.Fatalf("unexpected routes count: %d", len(resp.Routes))
	}

	if resp.Routes[0].Destination != "13.418555,52.521219" {
		t.Fatalf("unexpected first route: %+v", resp.Routes[0])
	}

	if resp.Routes[1].Destination != "13.428555,52.523219" {
		t.Fatalf("unexpected second route: %+v", resp.Routes[1])
	}

	if resp.Routes[2].Destination != "13.397634,52.529407" {
		t.Fatalf("unexpected third route: %+v", resp.Routes[2])
	}
}
