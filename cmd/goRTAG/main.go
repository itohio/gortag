package main

import "github.com/itohio/gortag/pkg/app"

func main() {
	rtag := app.New("Real Time Augio Generator")
	rtag.Run()
}
