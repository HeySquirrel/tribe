package main

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

func updateChanges(a *app.App) error {
	changed, err := changes()
	if err != nil {
		return err
	}

	a.UpdateChanges(changed)

	return nil
}

func update(a *app.App) {
	updateChanges(a)
	for {
		select {
		case <-a.Done:
			return
		case <-time.After(10 * time.Second):
			updateChanges(a)
		}
	}

}
