package data

import (
	"github.com/egoist/mygo/ui"

	"github.com/HycJack/MujicaUI/core"
	"github.com/HycJack/MujicaUI/internal"
	"github.com/HycJack/MujicaUI/theme"
)

// CardOptions configure a Card.
type CardOptions struct {
	// Title heads the card in the serif face; empty shows no header.
	Title string
	// Footer builds the bottom slot, below a divider; nil shows none.
	Footer func()
	// Decorated draws a double border and corner ornaments in Ornament.
	Decorated bool
	// HoverLift lifts the card and scales it up slightly while the pointer
	// is over it, as SwiftUI's cards; it takes no focus, and stays where it
	// is in the layout, animated by AnimatePose. The lift only affects the
	// way the card shows, never the space it reserves.
	HoverLift bool
}

// cardInset is how far the inner line of a decorated card sits inside.
const cardInset float32 = 4

// Card groups content on a surface with an outline. It adds no click behavior.
func Card(c *ui.Context, opts CardOptions, body func()) ui.Element {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	pad := u * 4
	if opts.Decorated {
		pad += cardInset + 2
	}
	card := ui.Column(c).Background(k.Surface).Radius(theme.CardRadius).Gap(u * 3).Padding(pad).Role(ui.RoleGroup)
	if opts.Title != "" {
		card.Label(opts.Title)
	}
	if opts.Decorated {
		card.Border(theme.BorderWidth, k.Ornament).DrawOver(func(p *ui.Painter, r ui.Rect) {
			internal.Frame(p, r, cardInset, 9, k.Ornament)
		})
	} else {
		card.Border(theme.BorderWidth, k.Border)
	}
	card.Children(func() {
		if opts.Title != "" {
			ui.Text(c, opts.Title).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size)).FontWeight(600).LineHeight(1.3)
			if opts.Decorated {
				ui.Box(c).Height(7).FillWidth().Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
					internal.RuleDiamond(p, r, 6, 7, k.Ornament)
				})
			}
		}
		ui.Column(c).Gap(u * 2).Children(body)
		if opts.Footer != nil {
			ui.Box(c).Height(1).FillWidth().Shrink(0).Background(k.Border)
			ui.Row(c).Gap(u * 2).Justify(ui.End).AlignItems(ui.Center).Wrap().Children(opts.Footer)
		}
	})
	cardPose(c, card, opts.HoverLift)
	return card
}

// cardLift is how far and how much a HoverLift card moves and scales.
const cardLift float32 = 4

// cardPose sets the hover lift of a Card, when lifted is true: it raises
// the card cardLift DIPs and scales it a hair while the pointer is over it,
// smoothly (AnimatePose), and back when it leaves. The lift only changes
// how the card shows, never the space it reserves in the layout. Without
// motion on the desktop the card shows there at once.
func cardPose(c *ui.Context, card ui.Element, lifted bool) {
	if !lifted {
		return
	}
	p := ui.Pose{}
	if card.Hovered() {
		p.Y = -cardLift
		p.Scale = 1.02
	}
	card.AnimatePose(p, ui.Animation{Duration: core.Motion(c, theme.StateDuration), Ease: ui.EaseInOut})
}
