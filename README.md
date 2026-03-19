# Tea Blends API

A simple API for managing tea blends.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

## Commands

The project uses `just` for command execution.

Build the application
```bash
just build
```

Run the application
```bash
just run
```

Create DB and other services container
```bash
just docker-run
```

Shutdown Container
```bash
just docker-down
```

Integrations Test:
```bash
just itest
```

Live reload the application:
```bash
just watch
```

Run the test suite:
```bash
just test
```

Clean up binary from the last build:
```bash
just clean
```
