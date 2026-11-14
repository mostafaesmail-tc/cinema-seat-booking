# Cinema Seat Booking

A lightweight cinema seat-booking demo built with Go and Redis. Browse a small
movie catalogue, view seat availability, temporarily hold a seat, then confirm
or release the reservation.

## Features

- Browser-based seat selection and checkout UI
- Redis-backed seat holds with a two-minute expiration
- Seat confirmation that persists the reservation in Redis
- Seat release and live availability polling in the UI
- Atomic hold creation so only one request can hold a seat at a time
- Optional Redis Commander UI for inspecting Redis data

## Tech stack

- Go 1.25.3+
- Redis 7
- Docker Compose for local Redis services

## Getting started

### Prerequisites

- Go 1.25.3 or later
- Docker and Docker Compose

### Run the application

1. Start Redis:

   ```sh
   docker compose up -d redis
   ```

2. Start the Go server from the repository root:

   ```sh
   go run ./cmd
   ```

3. Open [http://localhost:8080](http://localhost:8080).

The server connects to Redis at `localhost:6379` and serves the UI from the
`static/` directory. The optional Redis Commander service is available at
[http://localhost:8081](http://localhost:8081) when started with
`docker compose up -d`.

To stop the local services:

```sh
docker compose down
```

## Available movies

The current catalogue is defined in `cmd/main.go`:

| ID | Title | Seats |
|---|---|---|
| `inception` | Inception | 5 rows × 8 seats |
| `dune` | Dune: Part Two | 4 rows × 6 seats |

Seat IDs use a row letter and seat number, such as `A1` or `C4`.

## API

All endpoints are served from `http://localhost:8080`. JSON request bodies use
`user_id`; the current implementation treats it as a demo client identifier,
not as authentication.

### List movies

```http
GET /movies
```

Example response:

```json
[
  {
    "id": "inception",
    "title": "Inception",
    "rows": 5,
    "seats_per_row": 8
  }
]
```

### List booked or held seats

```http
GET /movies/{movieID}/seats
```

Returns an array of occupied seats. A `confirmed` value of `false` indicates a
temporary hold.

```json
[
  {
    "seat_id": "A1",
    "user_id": "demo-user",
    "booked": true,
    "confirmed": false
  }
]
```

### Hold a seat

```http
POST /movies/{movieID}/seats/{seatID}/hold
Content-Type: application/json

{
  "user_id": "demo-user"
}
```

A successful hold returns a session ID and expiration time. Holds expire after
two minutes unless confirmed.

### Confirm a hold

```http
PUT /sessions/{sessionID}/confirm
Content-Type: application/json

{
  "user_id": "demo-user"
}
```

Confirmation converts the temporary hold into a persistent booking.

### Release a hold

```http
DELETE /sessions/{sessionID}
Content-Type: application/json

{
  "user_id": "demo-user"
}
```

## Development

Run the Go test suite with Redis available at `localhost:6379`:

```sh
go test ./...
```

The booking concurrency test exercises many simultaneous attempts to hold the
same seat and requires a running Redis instance.

## Project structure

```text
cmd/                    Application entry point and movie catalogue
internal/adapters/redis/ Redis client setup
internal/booking/       Booking domain, service, HTTP handlers, and stores
internal/utils/         HTTP response helpers
static/                 Browser UI
docker-compose.yaml     Local Redis and Redis Commander services
```
