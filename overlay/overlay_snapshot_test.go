package overlay

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// overlayUntil advances the tester clock until cond holds, for at most 5 s of it.
func overlayUntil(tt *ui.Tester, cond func() bool) bool {
	for range 100 {
		if cond() {
			return true
		}
		tt.Advance(50 * time.Millisecond)
	}
	return cond()
}
