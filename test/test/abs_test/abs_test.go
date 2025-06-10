package abs_test

import (
	"testing"

	abs "github.com/LiangNing7/go-example/test/test"
)

func TestAbs(t *testing.T) {
	got := abs.Abs(-1)
	if got != 1 {
		t.Errorf("Abs(-1) = %f; want 1", got)
	}
}
