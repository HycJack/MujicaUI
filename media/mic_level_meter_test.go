package media

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/HycJack/MujicaUI/core"
)

func TestMicLevelMeterPeakHold(t *testing.T) {
	level := float32(1)
	var peak *micPeak
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		e := MicLevelMeter(c, level, MicLevelMeterOptions{Hold: 40 * time.Millisecond})
		peak = ui.Local(e, micPeak{}, func() micPeak { return micPeak{} })
	}, 400, 100)
	if !tt.HasText("0 dB") {
		t.Fatalf("texts %v", tt.Texts())
	}
	level = 0.1
	tt.Frame()
	if !tt.HasText("-20 dB") {
		t.Fatalf("level readout %v", tt.Texts())
	}
	if r, ok := tt.Find("Microphone level"); !ok || r.W == 0 {
		t.Fatal("meter not named")
	}
	if got := "peak " + micDB(peak.level); got != "peak 0 dB" {
		t.Fatalf("peak not held: %q", got)
	}
	time.Sleep(60 * time.Millisecond)
	tt.Frame()
	if got := "peak " + micDB(peak.level); got != "peak -20 dB" {
		t.Fatalf("peak did not fall after the hold: %q", got)
	}
}

func TestMicLevelMeterDB(t *testing.T) {
	for v, want := range map[float32]string{0: "-∞ dB", 1: "0 dB", 0.5: "-6 dB"} {
		if got := micDB(v); got != want {
			t.Errorf("%v: %s, want %s", v, got, want)
		}
	}
	mediaPanics(t, "over full scale", func(c *ui.Context) { MicLevelMeter(c, 1.2, MicLevelMeterOptions{}) })
}

func TestMicLevelMeterSnapshot(t *testing.T) {
	snapshot(t, "mic_level_meter", 360, 150, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(8).Children(func() {
			MicLevelMeter(c, 0.35, MicLevelMeterOptions{})
			MicLevelMeter(c, 0.8, MicLevelMeterOptions{})
			MicLevelMeter(c, 0.97, MicLevelMeterOptions{})
		})
	})
}
