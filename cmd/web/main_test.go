package main

import (
	"testing"
)

func TestRun(t *testing.T) {
	err := initializeAppConfig()
	if err != nil {
		t.Fatal("Failed initializeAppConfig()", err)
	}
	_, err = run()
	if err != nil {
		t.Fatal("Failed run()", err)
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
