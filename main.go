package main

import "time"
import "github.com/heysquirrel/tribe/app"
import "github.com/heysquirrel/tribe/git"

func main() {
	a := app.New()
	defer a.Close()

	a.UpdateChanges(git.Changes())
	// go update(a)

	a.Loop()
}

func update(a *app.App) {
	for {
		select {
		case <-a.Done:
			return
		case <-time.After(10 * time.Second):
			a.UpdateChanges(git.Changes())
		}
	}

}
