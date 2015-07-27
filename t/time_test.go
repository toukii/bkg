package t

import (
	"testing"
)

func TestNow(t *testing.T) {
	now := Now()
	t.Log(now)
}
