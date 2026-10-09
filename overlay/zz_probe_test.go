package overlay

import (
	"log/slog"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/HycJack/MujicaUI/core"
	"github.com/HycJack/MujicaUI/input"
)

func TestZZAlertProbe(t *testing.T) {
	slog.SetLogLoggerLevel(-4)
	open := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if b := input.Button(c, "Remove", input.ButtonOptions{}); b.Clicked() {
			open = true
		}
		AlertDialog(c, &open, AlertDialogOptions{Title: "Delete?", Description: "gone"})
	}, 600, 400)
	_ = tt.Click("Remove")
	tt.Key(0, ui.KeyEnter)
	t.Logf("after Enter: open=%v focusRemove=%v", open, tt.Focused("Remove"))
}
