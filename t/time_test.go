package t

import (
	"testing"
)

func TestLocNow(t *testing.T) {
	local := LocNow("Local")
	t.Logf("%s : %s \n", local.Location(), local.String())

	sh := LocNow("Asia/Shanghai")
	t.Logf("%s : %s \n", sh.Location(), sh.String())

	newYork := LocNow("America/New_York")
	t.Logf("%s : %s \n", newYork.Location(), newYork.String())
}
