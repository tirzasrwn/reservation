package main

import "testing"

func TestRun(t *testing.T) {
	_, err := run()
	if err != nil {
		t.Error("Failed run()")
	}
}

/**
run test:
1. Go to this folder.
2. Run $ go test -v
3. Run $ go test -cover
3. Run $ go test -coverprofile=coverage.out && go tool cover -html=coverage.out
4. See the result.
*/
