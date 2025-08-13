package main

import (
	"ping/config"
	"ping/tui"
)

func main() {
	config := config.Parse("config.yaml")

	tui.Start(config)
}
