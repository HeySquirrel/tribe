package main

import "github.com/heysquirrel/tribe/app"

func main() {
	a := app.New()
	defer a.Close()

	a.Loop()
}
