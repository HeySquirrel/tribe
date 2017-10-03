package widgets

import (
	"errors"
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

type keyBinding struct {
	view    string
	key     interface{}
	handler func()
}

// Focusable allows callers to work with ui widgets that can be focused
type Focusable interface {
	// Focus ensures the ui is updated to show that this widget is focused.
	// Focused widgets will also receive key events. If a widget isn't focused local key event handlers will not fire.
	Focus()

	// CanFocus will return true when a widget can be focused, false otherwise.
	CanFocus() bool

	// IsFocused will return true if the widget is currently focused, false otherwise.
	IsFocused() bool
}

// UI is the base cui widget that all other ui widgets should extend
type UI struct {
	Name    string
	Startx  float64
	Starty  float64
	Endx    float64
	Endy    float64
	Gui     *gocui.Gui
	FocusOn interface{}
	keys    []keyBinding
}

// Update this ui widget in a goroutine safe way. If the ui widget is updated with this
// method it is guarenteed to safe.
func (u *UI) Update(f func(v *gocui.View)) {
	u.Gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(u.Name)
		if err != nil {
			return err
		}

		f(v)
		return nil
	})
}

// Title changes the title of this ui widget
func (u *UI) Title(title string) {
	u.Update(func(v *gocui.View) {
		if u.FocusOn != nil {
			key, err := ToKeyString(u.FocusOn)
			if err != nil {
				log.Panicln(err)
			}
			v.Title = fmt.Sprintf("(%s) %s", key, title)
		} else {
			v.Title = title
		}
	})
}

// Focus ensures the ui is updated to show that this widget is focused.
// Focused widgets will also receive key events. If a widget isn't focused local key event handlers will not fire.
func (u *UI) Focus() {
	u.Update(func(v *gocui.View) {
		if u.Gui.CurrentView() != nil {
			u.Gui.CurrentView().Highlight = false
		}
		v.Highlight = true
		u.Gui.SetCurrentView(u.Name)
	})
}

// CanFocus will return true when a widget can be focused, false otherwise.
func (u *UI) CanFocus() bool {
	return u.FocusOn != nil
}

// IsFocused will return true if the widget is currently focused, false otherwise.
func (u *UI) IsFocused() bool {
	return u.Name == u.Gui.CurrentView().Name()
}

// Show ensures this widget is visible with respect to other widgets. Show will also make sure this widget
// is focused
func (u *UI) Show() func() {
	previousView := u.Gui.CurrentView()

	u.Update(func(v *gocui.View) {
		if previousView != nil {
			previousView.Highlight = false
		}
		u.Gui.SetCurrentView(u.Name)
		u.Gui.SetViewOnTop(u.Name)
	})

	return func() {
		u.Update(func(v *gocui.View) {
			v.Highlight = false
			u.Gui.SetViewOnBottom(u.Name)

			u.Gui.SetCurrentView(previousView.Name())
			previousView.Highlight = true
		})
	}
}

// AddLocalKey will add a key binding for this widget only. No other widgets will respond to
// this key's event handler
func (u *UI) AddLocalKey(key interface{}, handler func()) {
	u.keys = append(u.keys, keyBinding{u.Name, key, handler})
}

// AddGlobalKey will add a key binding for the entire app. It will respond no matter which
// widget is focused
func (u *UI) AddGlobalKey(key interface{}, handler func()) {
	u.keys = append(u.keys, keyBinding{"", key, handler})
}

// AddOneUseGlobalKey will add a global key binding that will only work the first time the key is pressed. As soon as the
// key handler executes the key binding is removed.
func (u *UI) AddOneUseGlobalKey(key interface{}, handler func()) {
	u.Update(func(v *gocui.View) {
		oneUseHandler := func() {
			u.Gui.DeleteKeybinding("", key, gocui.ModNone)
			handler()
		}
		u.Gui.SetKeybinding("", key, gocui.ModNone, ToBinding(oneUseHandler))
	})
}

// Layout positions this widget in the cui based on the Startx, Starty, Endx and Endy. It also registers
// all the key bindings that have been added to this widget
func (u UI) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	x1 := int(u.Startx * float64(maxX))
	y1 := int(u.Starty * float64(maxY))
	x2 := int(u.Endx*float64(maxX)) - 1
	y2 := int(u.Endy*float64(maxY)) - 1

	_, err := g.SetView(u.Name, x1, y1, x2, y2)
	if err != gocui.ErrUnknownView {
		return err
	}

	return u.registerKeyBindings(g)
}

func (u *UI) registerKeyBindings(g *gocui.Gui) error {
	if u.FocusOn != nil {
		u.AddGlobalKey(u.FocusOn, u.Focus)
	}

	for _, binding := range u.keys {
		err := g.SetKeybinding(binding.view, binding.key, gocui.ModNone, ToBinding(binding.handler))
		if err != nil {
			return err
		}
	}
	return nil
}

// ToBinding is a convenience method used to create keybinding methods from plain methods.
// This allows users to ignore details of the underlying cui framework
func ToBinding(f func()) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error { f(); return nil }
}

// ToKeyString converts a key to it's string representation
func ToKeyString(key interface{}) (string, error) {
	switch t := key.(type) {
	case gocui.Key:
		switch t {
		case gocui.KeyF1:
			return "F1", nil
		case gocui.KeyF2:
			return "F2", nil
		case gocui.KeyF3:
			return "F3", nil
		default:
			return "", errors.New("unknown key")
		}
	case rune:
		return string(t), nil
	default:
		return "", errors.New("unknown key")
	}
}
