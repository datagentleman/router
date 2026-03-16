package routes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	errInvalidResponse = errors.New("osrm: invalid response")
	errInvalidMetrics  = errors.New("osrm: invalid metrics")
)

const osrmBaseURL = "http://router.project-osrm.org"
const osrmTablePath = "/table/v1/driving/"

type osrmResponse struct {
	Code      string      `json:"code"`
	Durations [][]float64 `json:"durations"`
	Distances [][]float64 `json:"distances"`
}

// OSRMRouter implements Router interface.
type OSRMRouter struct {
	client *http.Client
}

// NewOSRMRouter creates an OSRM-backed router with a reusable HTTP client.
func NewOSRMRouter() OSRMRouter {
	router := OSRMRouter{
		client: &http.Client{Timeout: 5 * time.Second},
	}

	return router
}

// Fetch calls the OSRM table API and maps the response into internal metrics.
func (osrm OSRMRouter) Fetch(ctx context.Context, src string, dst []string) ([]Metrics, error) {
	url := osrm.buildURL(src, dst)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := osrm.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return osrm.parseResponse(resp, len(dst))
}

// buildURL creates the OSRM table request URL for one source and multiple destinations.
func (osrm OSRMRouter) buildURL(src string, dst []string) string {
	parts := make([]string, 0, len(dst)+1)
	parts = append(parts, src)
	parts = append(parts, dst...)

	destinations := make([]string, 0, len(dst))

	for i := range dst {
		destinations = append(destinations, strconv.Itoa(i+1))
	}

	values := url.Values{}
	values.Set("sources", "0")
	values.Set("destinations", strings.Join(destinations, ";"))
	values.Set("annotations", "distance,duration")

	coords := strings.Join(parts, ";")
	return osrmBaseURL + osrmTablePath + coords + "?" + values.Encode()
}

// parseResponse decodes the OSRM response and converts it into internal metrics.
func (osrm OSRMRouter) parseResponse(resp *http.Response, dstCount int) ([]Metrics, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, errInvalidResponse
	}

	var table osrmResponse
	if err := json.NewDecoder(resp.Body).Decode(&table); err != nil {
		return nil, errInvalidResponse
	}

	if table.Code != "Ok" {
		return nil, errInvalidMetrics
	}

	// We already checked for invalid response and !OK code, but we dont have any guarantee
	// that osrm returns valid durations/distances, so this is additional check.
	if len(table.Distances) == 0 || len(table.Durations) == 0 ||
		len(table.Distances[0]) < dstCount || len(table.Durations[0]) < dstCount {
		return nil, errInvalidMetrics
	}

	metrics := make([]Metrics, 0, dstCount)

	for i := 0; i < dstCount; i++ {
		metrics = append(metrics, Metrics{
			Distance: table.Distances[0][i],
			Duration: table.Durations[0][i],
		})
	}

	return metrics, nil
}
