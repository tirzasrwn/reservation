# Migration

## Prerequisites

- [golang-migrate](https://github.com/golang-migrate/migrate)

```sh
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

- postgresql

## Usage

```sh
# Export your postgres url:
export POSTGRESQL='postgres://user:password@host:port/dbname?sslmode=disable'
# Example:
export POSTGRESQL='postgres://postgres:postgres@localhost:5432/reservation?sslmode=disable'
# Create migration step:
migrate create -ext sql -dir migrations/ -seq create_user_table
# Migrate up:
migrate -database ${POSTGRESQL} -path migrations/ up
# Migrate down:
migrate -database ${POSTGRESQL} -path migrations/ down
```
