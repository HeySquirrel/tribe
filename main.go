package main

import "log"
import "github.com/jroimartin/gocui"
import "fmt"
import "os/exec"
import "time"
import "strings"
import "github.com/heysquirrel/tribe/app"

var (
	done = make(chan struct{})
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(app.Layout)

	err = keybindings(g)
	if err != nil {
		log.Panicln(err)
	}

	go update(g)

	err = g.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

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
	err := g.SetKeybinding("changes", gocui.KeyArrowDown, gocui.ModNone, cursorDown)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("changes", gocui.KeyArrowUp, gocui.ModNone, cursorUp)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	if err != nil {
		log.Panicln(err)
	}
	return nil
}

func changes() ([]string, error) {
	var results = make([]string, 1)

	cmdOut, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return nil, err
	}

	output := strings.Split(string(cmdOut), "\n")
	for _, change := range output {
		if len(change) > 0 {
			results = append(results, change[3:len(change)])
		}
	}

	return results, nil
}

func updateChanges(g *gocui.Gui) error {
	changed, err := changes()
	if err != nil {
		return err
	}

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("changes")
		if err != nil {
			return nil
		}
		v.Clear()
		for _, change := range changed {
			fmt.Fprintln(v, change)
		}
		return nil
	})

	return nil
}

func update(g *gocui.Gui) {
	updateChanges(g)
	for {
		select {
		case <-done:
			return
		case <-time.After(10 * time.Second):
			updateChanges(g)
		}
	}

}

func quit(g *gocui.Gui, v *gocui.View) error {
	close(done)
	return gocui.ErrQuit
}
