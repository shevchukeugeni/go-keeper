package main

import (
	"fmt"

	"keeper-project/cmd/client/app"
)

var (
	buildVersion = "N/A"
	buildTime    = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\n", buildVersion, buildTime)
	app.Execute()
}
