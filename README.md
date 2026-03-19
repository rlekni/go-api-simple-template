# Tea Blends API

A simple API for managing tea blends.

## Bill of Materials (BOM)

- **Language:** [Go 1.26+](https://go.dev/)
- **API Framework:** [Huma v2](https://huma.rocks/)
- **Router:** [Chi v5](https://github.com/go-chi/chi)
- **Database:** [PostgreSQL](https://www.postgresql.org/)
- **Database Driver:** [PGX v5](https://github.com/jackc/pgx)
- **SQL Generator:** [SQLC](https://sqlc.dev/)
- **Migrations:** [Goose](https://github.com/pressly/goose)
- **Observability:**
  - [OpenTelemetry](https://opentelemetry.io/) (Tracing, Metrics, Logs)
  - [Statsviz](https://github.com/arl/statsviz) (Runtime metrics visualization)
- **Tooling:**
  - [Just](https://github.com/casey/just) (Command runner)
  - [Air](https://github.com/air-verse/air) (Live reload)
  - [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

## Commands

The project uses `just` for command execution.

### SQL Generation
Regenerate Go code from SQL queries (cleans `internal/postgresql` first):
```bash
just generate
```

### Build & Run
Build the application:
```bash
just build
```

Run the application:
```bash
just run
```

Live reload the application (requires `air`):
```bash
just watch
```

### Infrastructure
Create DB and OTel collector containers:
```bash
just docker-run
```

Shutdown containers:
```bash
just docker-down
```

### Testing & Quality
Run the test suite:
```bash
just test
```

Integrations Test:
```bash
just itest
```

### Cleanup
Clean up binary and generated docs:
```bash
just clean
```
