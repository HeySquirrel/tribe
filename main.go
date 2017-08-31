package main

import "time"
import "github.com/heysquirrel/tribe/app"

func main() {
	a := app.New()
	defer a.Close()

	go update(a)

	a.Loop()
}

func update(a *app.App) {
	a.UpdateChanges(a.Git.Changes())
	for {
		select {
		case <-a.Done:
			return
		case <-time.After(10 * time.Second):
			a.Debug("Checking for changes")
			a.UpdateChanges(a.Git.Changes())
		}
	}

}
