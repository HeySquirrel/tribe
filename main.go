package main

import "github.com/jroimartin/gocui"
import "fmt"
import "os/exec"
import "time"
import "strings"
import "github.com/heysquirrel/tribe/app"

func main() {
	a := app.New()
	defer a.Close()

	go update(a)

	a.Loop()
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

func updateView(g *gocui.Gui, view string, value string) {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(view)
		if err != nil {
			return nil
		}
		v.Clear()
		fmt.Fprintln(v, value)
		return nil
	})
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

func update(a *app.App) {
	updateChanges(a.Gui)
	for {
		select {
		case <-a.Done:
			return
		case <-time.After(10 * time.Second):
			updateChanges(a.Gui)
		}
	}

}
