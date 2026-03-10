package main

import (
	"os"

	"github.com/bluefunda/abaper-cli/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
