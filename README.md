# reservation

## Setup
```sh
go mod tidy
```

## Running

### Web Application

```sh
go run $(ls cmd/web/*.go | grep -v _test.go)
```
### Test

```sh
# Go to inside the directory part you want to test
go test -v
```
Example
```sh
# From root level project, go to interna/handlers/
go test -v
```