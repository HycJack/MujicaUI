package overlay

import (
	"log/slog"

	"github.com/egoist/mygo/ui"
)

// ContextMenuOptions configure a ContextMenu.
type ContextMenuOptions struct {
	Disabled bool
}

// contextMenuState is where the menu opened, relative to the target.
type contextMenuState struct {
	open     bool
	x, y     float32
	hadFocus bool // the focus was inside the menu last frame
}

// ContextMenu gives the container target an in-window menu at the pointer or below the focus.
func ContextMenu(c *ui.Context, target ui.Element, opts ContextMenuOptions, build func(m *PopupMenu)) ui.Element {
	if !target.Valid() {
		panic("mujica: ContextMenu needs a target")
	}
	if build == nil {
		panic("mujica: ContextMenu needs a build function")
	}
	st := ui.Local(target, "mujica.contextmenu", func() contextMenuState { return contextMenuState{} })
	if opts.Disabled {
		st.open = false
		return target
	}
	target.HandleInput(func(ev ui.InputEvent) bool {
		// Take the press so the release comes back, then open on it: the
		// press itself would read as a press outside the panel and close
		// the menu in the same frame.
		switch ev.Kind {
		case ui.InputPointerDown:
			return ev.Button == 1
		case ui.InputPointerUp:
			if ev.Button != 1 {
				return false
			}
			st.open, st.x, st.y = true, ev.X, ev.Y
			slog.Debug("mujica: context_menu open", "by", "pointer")
			return true
		}
		return false
	})
	if target.Shortcut(0, ui.KeyContextMenu) || target.Shortcut(ui.Shift, ui.KeyF10) {
		st.open, st.x, st.y = true, 0, target.Bounds().H
		slog.Debug("mujica: context_menu open", "by", "keyboard")
	}
	if !st.open {
		return target
	}
	var anchor ui.Element
	target.Children(func() {
		anchor = ui.Box(c).Key("mujica.contextmenu.anchor").Absolute().Left(st.x).Top(st.y).Size(1, 1).PassThrough()
	})
	panel := PopupMenuOpen(c, anchor, &st.open, build)
	focusIn := panel.Valid() && panel.FocusWithin()
	if !st.open {
		// Escape or a choice: give the target the focus the menu held.
		if st.hadFocus || focusIn {
			target.Focus()
		}
		st.hadFocus = false
	} else {
		st.hadFocus = focusIn
	}
	return target
}
