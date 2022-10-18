package main

import "ewintr.nl/gte/cmd/android-app/runner"

func main() {
	runner := runner.NewRunner()
	runner.Init()
	runner.Run()
}
