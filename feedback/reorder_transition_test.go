package feedback

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/HycJack/MujicaUI/core"
	"github.com/HycJack/MujicaUI/theme"
)

func TestReorderTransitionSlides(t *testing.T) {
	keys := []string{"a", "b", "c"}
	tt, n := motionFrames(func(c *ui.Context) {
		ReorderTransition(c, keys, ReorderTransitionOptions{Gap: 4}, func(k string) {
			ui.Text(c, "Item "+k).Padding(8)
		})
	}, 300, 200)
	ra, _ := tt.Find("Item a")
	keys = []string{"c", "b", "a"}
	before := *n
	tt.Frame()
	if *n-before < 3 {
		t.Fatal("reorder did not animate")
	}
	// Half a move in, a sits between where it was and where it lands.
	motionAdvance(tt, theme.StateDuration/2)
	mid, _ := tt.Find("Item a")
	motionSettle(tt, theme.StateDuration)
	end, _ := tt.Find("Item a")
	if !(mid.Y > ra.Y && mid.Y < end.Y) || motionBusy(tt, n) {
		t.Fatalf("a moved %v → %v → %v", ra.Y, mid.Y, end.Y)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("duplicate key did not panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ReorderTransition(c, []int{1, 1}, ReorderTransitionOptions{}, func(int) {})
	}, 100, 50, 1)
}

func TestReorderTransitionSnapshot(t *testing.T) {
	snapshot(t, "reorder_transition", 300, 140, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			ReorderTransition(c, []string{"Sakiko", "Uika", "Mutsumi"}, ReorderTransitionOptions{Gap: 4}, func(k string) {
				ui.Text(c, k).Padding(4, 8).Background(core.Tokens(c).Surface).Radius(4)
			})
		})
	})
}

func TestReorderTransitionInvalidGap(t *testing.T) {
	for _, gap := range []float32{float32(nanValue()), float32(math.Inf(1))} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("invalid gap %v accepted", gap)
				}
			}()
			ui.Render(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				ReorderTransition(c, []int{1}, ReorderTransitionOptions{Gap: gap}, func(int) {})
			}, 200, 100, 1)
		}()
	}
}

func TestReorderTransitionReducedMotionCancelsExit(t *testing.T) {
	keys := []string{"a", "b"}
	tt, n := motionFrames(func(c *ui.Context) {
		ReorderTransition(c, keys, ReorderTransitionOptions{Duration: 5 * time.Second}, func(k string) { ui.Text(c, k) })
	}, 300, 150)
	keys = []string{"b"}
	tt.Frame()
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if motionBusy(tt, n) || tt.HasText("a") {
		t.Fatal("removed item kept animating under reduced motion")
	}
}
