package pow_test

import (
	"testing"

	"github.com/cdonut/powddos/pkg/pow"
)

var attempts uint64 = 1 << 30

func TestPowSolve(t *testing.T) {
	c := pow.NewChallenge(2, "test")

	ans, err := c.Solve(attempts)

	if err != nil {
		t.Errorf("pow solution failed: %s", err.Error())
	}

	if pow.CheckSolution(ans, 2, "test", 300) != true {
		t.Error("pow check failed!")
	}
}
