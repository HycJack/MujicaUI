package finance

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/HycJack/MujicaUI/data"
)

func TestSentimentGaugeBandsAndBoundary(t *testing.T) {
	for i, value := range []float64{0, 20, 40, 60, 80, 100} {
		tt := chatTester(400, 300, func(c *ui.Context) { SentimentGauge(c, value, SentimentGaugeOptions{Label: "Sentiment"}) })
		words := []string{"Extreme fear", "Fear", "Neutral", "Greed", "Extreme greed", "Extreme greed"}
		if !tt.HasText(words[i]) {
			t.Fatalf("value %v texts %q", value, tt.Texts())
		}
	}
	marketsPanic(t, func(c *ui.Context) { SentimentGauge(c, 101, SentimentGaugeOptions{}) })
}
func TestSentimentGaugeSnapshot(t *testing.T) {

	snapshot(t, "sentiment_gauge", 320, 280, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Fill().Children(func() {
			SentimentGauge(c, 72, SentimentGaugeOptions{Label: "Fear & Greed"})
		})
	})
}

func TestSentimentGaugeKeyboardReadout(t *testing.T) {
	tt := chatTester(400, 300, func(c *ui.Context) { SentimentGauge(c, 72, SentimentGaugeOptions{Label: "Sentiment"}) })
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if len(tt.Announcements()) == 0 {
		t.Fatal("missing keyboard readout")
	}
}

func TestSentimentGaugeStatesKeepSize(t *testing.T) {
	status := data.DataReady
	empty := false
	retried := 0
	var heights []float32
	tt := chatTester(400, 300, func(c *ui.Context) {
		element := SentimentGauge(c, 72, SentimentGaugeOptions{Label: "Sentiment", Status: status, Empty: empty, Error: "Caller failed", Retry: func() { retried++ }})
		if h := element.Bounds().H; h > 0 && (len(heights) == 0 || heights[len(heights)-1] != h) {
			heights = append(heights, h) // bounds read at build time, one sample per state
		}
	})
	status = data.DataLoading
	tt.Frame()
	status = data.DataFailed
	tt.Frame()
	_ = tt.Click("Retry")
	if retried != 1 {
		t.Fatalf("retries %d", retried)
	}
	status = data.DataReady
	empty = true
	tt.Frame()
	if !tt.HasText("No data") {
		t.Fatalf("empty state texts %q", tt.Texts())
	}
	for _, h := range heights {
		if h != heights[0] {
			t.Fatalf("heights changed across states: %v", heights)
		}
	}
	if len(heights) != 1 || heights[0] <= 0 {
		t.Fatalf("heights not sampled cleanly: %v", heights)
	}
}
