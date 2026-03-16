package routes

import "sort"

type Route struct {
	Destination string  `json:"destination"`
	Duration    float64 `json:"duration"`
	Distance    float64 `json:"distance"`
}

type Response struct {
	Source string  `json:"source"`
	Routes []Route `json:"routes"`
}

// newRouteResponse builds the outward API response and sorts routes by duration and distance.
func newRouteResponse(src string, dst []string, metrics []Metrics) Response {
	routes := make([]Route, 0, len(dst))

	for i, destination := range dst {
		routes = append(routes, Route{
			Destination: destination,
			Duration:    metrics[i].Duration,
			Distance:    metrics[i].Distance,
		})
	}

	sortRoutes(routes)

	return Response{
		Source: src,
		Routes: routes,
	}
}

// sortRoutes sorts routes by duration first and distance second.
func sortRoutes(routes []Route) {
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Duration == routes[j].Duration {
			return routes[i].Distance < routes[j].Distance
		}

		return routes[i].Duration < routes[j].Duration
	})
}
