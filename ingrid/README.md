# Router Routes Service

Small Go HTTP service that accepts one source coordinate and multiple destination coordinates, and returns routes sorted by duration and distance.

## Setup

Requirements:

- Go 1.22+

Check your Go version:

```sh
go version
```

Clone or copy the repository and move into the project root:

```sh
cd router
```

## Build

```sh
go build -o routes ./cmd/routes
```

## Run

```sh
./routes
```

The service listens on `:8080`.

Example request:

```sh
curl "http://localhost:8080/route?src=13.388860,52.517037&dst=13.397634,52.529407&dst=13.428555,52.523219"
```

## Test

Run the default test suite:

```sh
go test ./...
```

This runs the default unit test suite only.

Live OSRM integration tests are disabled by default. They call the real OSRM third-party API, so they require internet access and depend on the external service being available.

Run all tests including live OSRM integration tests:

```sh
INTEGRATION=1 go test ./...
```

Run a specific test:

```sh
INTEGRATION=1 go test -run TestOSRMRouterFetchSuccess ./...
```

## Tradeoffs

- The service does not implement rate limiting or throttling yet. This could be added in the application itself or at the edge through a reverse proxy / gateway.
- There is no caching yet, so every request currently depends on the third-party OSRM API.
- This keeps the implementation simple, but it also means overall latency and availability are influenced by the upstream service.
- In the future, responses from 3rd party APIs can be mocked to make tests faster.


## Scalability

What can be done to scale this service in the future

- Run multiple stateless instances of the service behind a load balancer.
- Add caching for `src -> dst` route metrics, for example in memory or in Redis, to reduce repeated calls to the upstream router.
- Reusing cached route metrics would allow the service to handle higher throughput with less dependency on the external API.

## System resilience

- The service is stateless, so restarts are simple and do not require recovery of local state.
- Running multiple instances behind a load balancer would improve availability if one instance becomes unhealthy.
- Monitoring and health checks could be added to detect failures earlier and make the service easier to operate.
