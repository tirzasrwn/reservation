# reservation

## Setup
```sh
go mod tidy
```

## Running

### Web Application

```sh
go run $(ls cmd/web/*.go | grep -v _test.go)
# Or just run run.sh script.
chmod +x run.sh
./run.sh
```
### Test

```sh
# Go to inside the directory part you want to test
go test -v
go test -cover # To look pecentage of covarage test.
go test -coverprofile=coverage.out && go tool cover -html=coverage.out # Same as cover but with more detail and html format.

```
Example
```sh
# From root level project, go to internal/handlers/
go test -v
```