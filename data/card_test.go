package data

import (
	"image/color"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/HycJack/MujicaUI/core"
	"github.com/HycJack/MujicaUI/input"
	"github.com/HycJack/MujicaUI/theme"
)

func cardView(decorated bool) func(c *ui.Context) {
	return func(c *ui.Context) {
		ui.Column(c).Padding(20).AlignItems(ui.Start).Children(func() {
			Card(c, CardOptions{Title: "Archive", Decorated: decorated, Footer: func() {
				input.Button(c, "Open", input.ButtonOptions{})
			}}, func() {
				ui.Text(c, "Letters of the court")
			}).Width(260)
		})
	}
}

func TestCardSlots(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		cardView(false)(c)
	}, 400, 300)
	for _, s := range []string{"Archive", "Letters of the court", "Open"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
	title, _ := tt.Find("Archive")
	body, _ := tt.Find("Letters of the court")
	foot, _ := tt.Find("Open")
	if !(title.Y < body.Y && body.Y < foot.Y) {
		t.Fatalf("slot order: title %v body %v footer %v", title.Y, body.Y, foot.Y)
	}
}

func TestCardNoTitleNoFooter(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Card(c, CardOptions{}, func() { ui.Text(c, "Only body") })
	}, 300, 200)
	if got := tt.Texts(); len(got) != 1 || got[0] != "Only body" {
		t.Fatalf("texts %q", got)
	}
}

func TestCardDecorationColor(t *testing.T) {
	for _, dec := range []bool{false, true} {
		var bounds ui.Rect
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{Mode: core.Light})
			ui.Column(c).Padding(20).AlignItems(ui.Start).Children(func() {
				bounds = Card(c, CardOptions{Title: "T", Decorated: dec}, func() {}).Width(200).Bounds()
			})
		}, 300, 200)
		tt.Frame()
		img := tt.Image()
		// The inner corner diamond sits at the inset corner.
		px := img.At(int(bounds.X+cardInset), int(bounds.Y+cardInset))
		want := toRGBA(theme.Light().Ornament)
		if got := color.RGBAModel.Convert(px).(color.RGBA); (got == want) != dec {
			t.Fatalf("decorated=%v: corner pixel %v, ornament %v", dec, got, want)
		}
	}
}

func toRGBA(c ui.Color) color.RGBA { return color.RGBA{c.R, c.G, c.B, c.A} }

// firstSurfaceY returns the first y where, at column x, the image shows the
// card's surface color, scanning down from top to bottom; -1 when the column
// never shows it. It skips the card's top radius and border by starting at
// the layout top, so it tracks where the filled body starts.
func firstSurfaceY(t *testing.T, tt *ui.Tester, surf color.RGBA, x, top, bottom int) int {
	t.Helper()
	for y := top; y < bottom; y++ {
		if got := color.RGBAModel.Convert(tt.Image().At(x, y)).(color.RGBA); got == surf {
			return y
		}
	}
	return -1
}

func TestCardHoverLiftMovesOnHover(t *testing.T) {
	var surf color.RGBA
	var laytop float32
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		k := core.Tokens(c)
		surf = color.RGBA{k.Surface.R, k.Surface.G, k.Surface.B, k.Surface.A}
		ui.Column(c).FillWidth().FillHeight().Background(k.Background).Padding(20).Children(func() {
			laytop = Card(c, CardOptions{HoverLift: true}, func() { ui.Text(c, "lift") }).Width(120).Bounds().Y
		})
	}, 220, 180)
	tt.Frame()
	// The probe runs down the card's middle, past its top radius.
	start := int(laytop) - 2
	hoverY := firstSurfaceY(t, tt, surf, 100, start, start+14)
	tt.Move(60, 40) // over the card
	tt.Advance(theme.StateDuration)
	liftY := firstSurfaceY(t, tt, surf, 100, start, start+14)
	tt.Move(5, 5) // off the card
	tt.Advance(theme.StateDuration)
	restY := firstSurfaceY(t, tt, surf, 100, start, start+14)
	if hoverY < 0 || liftY < 0 || restY < 0 {
		t.Fatalf("surface not found: at rest %d, hovering %d, back %d", hoverY, liftY, restY)
	}
	if liftY >= hoverY {
		t.Fatalf("hover lift: body top at rest %d, hovering %d; want it to rise", hoverY, liftY)
	}
	if restY != hoverY {
		t.Fatalf("back to rest: body top %d, want %d", restY, hoverY)
	}
	// The layout keeps its place: the lift only changes how it shows.
	tt.Frame()
	if b := firstSurfaceY(t, tt, surf, 100, int(laytop)-2, int(laytop)+6); b != hoverY {
		t.Fatalf("layout bounds unchanged, body still at %d, got %d", hoverY, b)
	}
}

func TestCardNoLiftWithoutOption(t *testing.T) {
	var surf color.RGBA
	var laytop float32
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		k := core.Tokens(c)
		surf = color.RGBA{k.Surface.R, k.Surface.G, k.Surface.B, k.Surface.A}
		ui.Column(c).FillWidth().FillHeight().Background(k.Background).Padding(20).Children(func() {
			laytop = Card(c, CardOptions{}, func() { ui.Text(c, "lift") }).Width(120).Bounds().Y
		})
	}, 220, 180)
	tt.Frame()
	start := int(laytop) - 2
	atRest := firstSurfaceY(t, tt, surf, 100, start, start+14)
	tt.Move(60, 40)
	tt.Advance(theme.StateDuration)
	hovering := firstSurfaceY(t, tt, surf, 100, start, start+14)
	if atRest < 0 || hovering < 0 {
		t.Fatalf("surface not found: rest %d, hover %d", atRest, hovering)
	}
	// Without HoverLift the card shows in place; hovering must not move it.
	if atRest != hovering {
		t.Fatalf("card without HoverLift moved: rest top %d, hovering top %d", atRest, hovering)
	}
}

func TestCardSnapshot(t *testing.T) {
	snapshot(t, "card", 640, 240, func(c *ui.Context) {
		ui.Row(c).Padding(20).Gap(20).AlignItems(ui.Start).Children(func() {
			Card(c, CardOptions{Title: "Archive", Footer: func() {
				input.Button(c, "Open", input.ButtonOptions{Variant: input.Secondary})
			}}, func() { ui.Text(c, "Letters of the court, sealed in wine-red wax.") }).Width(280)
			Card(c, CardOptions{Title: "Chapel", Decorated: true, Footer: func() {
				input.Button(c, "Enter", input.ButtonOptions{})
			}}, func() { ui.Text(c, "Vespers at seven; candles by the north door.") }).Width(280)
		})
	})
	for _, d := range []theme.Density{theme.Compact, theme.Comfortable} {
		snapshot(t, "card-"+strings.ToLower(d.String()), 640, 260, func(c *ui.Context) {
			mode := core.Light
			if core.Tokens(c).Dark {
				mode = core.Dark
			}
			core.Use(c, core.Settings{Mode: mode, Density: d})
			cardView(true)(c)
		})
	}
}
