// Package theme holds MujicaUI's design tokens: colors, type, metrics and motion.
package theme

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// Tokens are the colors of one appearance.
type Tokens struct {
	Dark          bool
	Background    ui.Color
	Surface       ui.Color
	SurfaceHover  ui.Color
	Text          ui.Color
	TextMuted     ui.Color
	Border        ui.Color
	ControlBorder ui.Color
	Accent        ui.Color
	AccentHover   ui.Color
	AccentPressed ui.Color
	OnAccent      ui.Color
	AccentText    ui.Color
	Selection     ui.Color
	Focus         ui.Color
	Ornament      ui.Color
	Success       ui.Color
	Warning       ui.Color
	Danger        ui.Color
	Info          ui.Color
	DangerFill    ui.Color
	OnDanger      ui.Color
}

// Light returns the Mint appearance: a zinc grayscale on white, with blue
// reserved for accent — the palette keeps color scarce and surfaces layered
// (white cards on grey panels), so hierarchy comes from the surface steps.
func Light() Tokens {
	return Tokens{
		Background:    ui.Hex("#FFFFFF"),
		Surface:       ui.Hex("#F4F4F5"),
		SurfaceHover:  ui.Hex("#ECECEE"),
		Text:          ui.Hex("#18181B"),
		TextMuted:     ui.Hex("#6B6B74"),
		Border:        ui.Hex("#D9D9DE"),
		ControlBorder: ui.Hex("#8B8B95"),
		Accent:        ui.Hex("#2563EB"),
		AccentHover:   ui.Hex("#1D4ED8"),
		AccentPressed: ui.Hex("#1E40AF"),
		OnAccent:      ui.Hex("#FFFFFF"),
		AccentText:    ui.Hex("#1D4ED8"),
		Selection:     ui.Hex("#E8EFFD"),
		Focus:         ui.Hex("#2563EB"),
		Ornament:      ui.Hex("#8B8B95"),
		Success:       ui.Hex("#1F7A3D"),
		Warning:       ui.Hex("#9A5410"),
		Danger:        ui.Hex("#C62B30"),
		Info:          ui.Hex("#0E7490"),
		DangerFill:    ui.Hex("#C62B30"),
		OnDanger:      ui.Hex("#FFFFFF"),
	}
}

// Dark returns the Mint appearance inverted: the same zinc steps on
// near-black, the accent lifted to hold contrast.
func Dark() Tokens {
	return Tokens{
		Dark:          true,
		Background:    ui.Hex("#111113"),
		Surface:       ui.Hex("#1C1C1F"),
		SurfaceHover:  ui.Hex("#26262A"),
		Text:          ui.Hex("#F4F4F5"),
		TextMuted:     ui.Hex("#A1A1AA"),
		Border:        ui.Hex("#33333A"),
		ControlBorder: ui.Hex("#6E6E78"),
		Accent:        ui.Hex("#5B8DFF"),
		AccentHover:   ui.Hex("#7AA5FF"),
		AccentPressed: ui.Hex("#4A7AF0"),
		OnAccent:      ui.Hex("#111113"),
		AccentText:    ui.Hex("#9DC0FF"),
		Selection:     ui.Hex("#16233D"),
		Focus:         ui.Hex("#9DC0FF"),
		Ornament:      ui.Hex("#71717A"),
		Success:       ui.Hex("#5FC47C"),
		Warning:       ui.Hex("#E5A04A"),
		Danger:        ui.Hex("#FF6B6F"),
		Info:          ui.Hex("#22D3EE"),
		DangerFill:    ui.Hex("#C62B30"),
		OnDanger:      ui.Hex("#FFFFFF"),
	}
}

// HighContrast returns t with secondary text and borders raised for readability.
func (t Tokens) HighContrast() Tokens {
	t.TextMuted = t.Text
	t.Border = t.ControlBorder
	return t
}

// Density sets control heights and the spacing unit.
type Density int

const (
	Default Density = iota
	Compact
	Comfortable
)

// Unit is the base spacing unit in DIPs.
func (d Density) Unit() float32 {
	switch d {
	case Compact:
		return 3
	case Comfortable:
		return 5
	}
	return 4
}

// ControlHeight is the minimum height of single-line controls and data rows.
func (d Density) ControlHeight() float32 {
	switch d {
	case Compact:
		return 28
	case Comfortable:
		return 44
	}
	return 36
}

func (d Density) String() string {
	switch d {
	case Compact:
		return "Compact"
	case Comfortable:
		return "Comfortable"
	}
	return "Default"
}

// Type sizes and families.
const (
	H1Size      float32 = 28
	H2Size      float32 = 22
	H3Size      float32 = 18
	BodySize    float32 = 14
	CaptionSize float32 = 12
	LineHeight  float32 = 1.5

	SerifFont = "Georgia, Songti SC, SimSun, Noto Serif CJK SC, serif"
	MonoFont  = "monospace"
)

// Shape metrics in DIPs. One radius family, so curves read as related;
// the values are generous — pills and switches go past them via Pill uses.
const (
	ControlRadius float32 = 14
	CardRadius    float32 = 18
	PopupRadius   float32 = 24
	BorderWidth   float32 = 1
	FocusWidth    float32 = 2
	FocusGap      float32 = 2
	MinHitSize    float32 = 28
)

// Motion durations.
const (
	StateDuration = 120 * time.Millisecond
	PopupDuration = 160 * time.Millisecond
)
