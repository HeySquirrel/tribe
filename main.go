package main

import "log"
import "github.com/jroimartin/gocui"
import "fmt"
import "os/exec"

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("side", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	return nil
}

func changes() (string, error) {
	var (
		cmdOut []byte
		err    error
	)

	if cmdOut, err = exec.Command("git", "status", "--porcelain").Output(); err != nil {
		return "", err
	}
	output := string(cmdOut)
	return output[3:len(output)], nil
}

func update(g *gocui.Gui) error {
	var (
		changed string
		err     error
	)

	if changed, err = changes(); err != nil {
		return err
	}

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("side")
		if err != nil {
			return nil
		}
		v.Clear()
		fmt.Fprint(v, changed)
		return nil
	})

	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	x1 := int(0.0 * float64(maxX))
	y1 := int(0.0 * float64(maxY))
	x2 := int(0.1*float64(maxX)) - 1
	y2 := int(0.3*float64(maxY)) - 1

	if v, err := g.SetView("side", x1, y1, x2, y2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Changes"
		v.Highlight = true
		v.SelBgColor = gocui.ColorCyan
		v.SelFgColor = gocui.ColorBlack
	}

	if _, err := g.SetCurrentView("side"); err != nil {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	update(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
