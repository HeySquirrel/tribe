package main

import "time"
import "github.com/heysquirrel/tribe/app"
import "github.com/heysquirrel/tribe/git"

func main() {
	a := app.New()
	defer a.Close()

	go update(a)

	a.Loop()
}

func update(a *app.App) {
	a.UpdateChanges(git.Changes())
	for {
		select {
		case <-a.Done:
			return
		case <-time.After(10 * time.Second):
			a.Debug("Checking for changes")
			a.UpdateChanges(git.Changes())
		}
	}

}
