package main

import (
	"fmt"
	"os"

	"github.com/FRAZ5094/ping/config"
	"github.com/FRAZ5094/ping/tui"
)

func main() {
	hosts, err := config.Parse("config.yaml")
	if err != nil {
		fmt.Println("Error parsing config:", err)
		os.Exit(1)
	}

	tui.Start(hosts)
}
