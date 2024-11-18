# ***kera*** - Habit Tracker

Fully functional **habit tracker** backend providing an API for managing users, sessions and habits.

Built with Go's standard library and PostgreSQL ([pgx](https://github.com/jackc/pgx) driver).

## Prerequisites

- [Go](https://go.dev) (developed and tested with Go 1.23)
- [PostgreSQL](https://www.postgresql.org/) database

## Usage

The app requires ADDRESS and DSN (Data Source Name) environment variables.

On Linux:
```bash
cd kera
export ADDRESS=address:port
export DSN=postgres://username:password@address:port/database
go run .
```

## API documentation

The OpenAPI specification file is available [here](./openapi.yaml).
You can visualize it using a tool like Swagger UI.

## Disclaimer

This is a personal project created for learning purposes and is **not suitable** for real-world usage.

## License

This project is licensed under [AGPL-3.0](./LICENSE).
