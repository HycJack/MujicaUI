package feedback

import (
	"log/slog"
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/HycJack/MujicaUI/core"
	"github.com/HycJack/MujicaUI/input"
)

// TypewriterOptions configure a Typewriter.
type TypewriterOptions struct {
	// Speed is characters per second; zero means 40.
	Speed float64
	// Skippable adds a Skip button.
	Skippable bool
}

type typewriterState struct {
	text    []rune
	shown   int
	last    time.Time
	skipped bool
	done    bool
}

// Typewriter reveals text one character at a time.
func Typewriter(c *ui.Context, text string, opts TypewriterOptions) core.ChoiceResult {
	speed := opts.Speed
	if speed == 0 {
		speed = 40
	}
	if math.IsNaN(speed) || math.IsInf(speed, 0) || speed < float64(time.Second)/float64(math.MaxInt64) || speed > float64(time.Second) {
		panic("mujica: Typewriter Speed must give a finite positive interval")
	}
	// The slowest speed asks for the longest interval, which rounds past the
	// largest Duration: it types one character per Duration and no faster.
	step := time.Duration(math.MaxInt64)
	if interval := float64(time.Second) / speed; interval < float64(math.MaxInt64) {
		step = time.Duration(interval)
	}
	col := ui.Column(c).AlignItems(ui.Start).Gap(core.Density(c).Unit() * 2).Role(ui.RoleGroup).Label(text)
	res := core.ChoiceResult{Element: col}
	st := ui.Local(col, "typing", func() *typewriterState { return &typewriterState{last: c.Now()} })
	s := *st
	runes := []rune(text)
	if !strings.HasPrefix(text, string(s.text)) {
		s.shown, s.skipped, s.done = 0, false, false
		s.last = c.Now()
		slog.Debug("mujica: typewriter restart")
	}
	if len(runes) > len(s.text) && s.done {
		// Appended text continues at the typing speed.
		s.done, s.skipped, s.last = false, false, c.Now()
	}
	s.text = runes
	now := c.Now()
	if c.Preferences().ReduceMotion || s.skipped {
		s.shown = len(runes)
	} else if s.shown < len(runes) {
		if n := int(now.Sub(s.last) / step); n > 0 {
			s.shown = min(len(runes), s.shown+n)
			s.last = s.last.Add(time.Duration(n) * step)
		}
		if s.shown < len(runes) {
			col.DrawOver(func(*ui.Painter, ui.Rect) { c.After(s.last.Add(step).Sub(c.Now())) })
		}
	} else {
		s.last = now
	}
	col.Children(func() {
		ui.Text(c, string(runes[:s.shown]))
		if opts.Skippable && s.shown < len(runes) {
			b := input.Button(c, core.Msg(c, "typewriter.skip"), input.ButtonOptions{Variant: input.Ghost})
			if b.Clicked() && !b.IsDisabled() {
				s.skipped, s.shown = true, len(runes)
				slog.Debug("mujica: typewriter skip")
			}
		}
	})
	if s.shown == len(runes) && !s.done {
		s.done, res = true, core.NewChoiceResult(res.Element, true)
	}
	return res
}
