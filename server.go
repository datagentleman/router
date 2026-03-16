package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type Metrics struct {
	Distance float64
	Duration float64
}

// Router fetches route metrics for one source and multiple destinations.
type Router interface {
	Fetch(ctx context.Context, src string, dst []string) ([]Metrics, error)
}

var router = NewOSRMRouter()

// Run starts the HTTP server.
func Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /route", routeHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("listening on http://localhost%s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

// routeHandler validates request params, fetches route metrics, and writes the API response.
func routeHandler(w http.ResponseWriter, r *http.Request) {
	req := parseRequest(r)

	err := req.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metrics, err := router.Fetch(r.Context(), req.Src, req.Dst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(newResponse(req.Src, req.Dst, metrics))
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
