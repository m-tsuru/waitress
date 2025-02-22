package main

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed VERSION
var version string

var usage = fmt.Sprintf(`
waitress %s
`, version)[1:]

func main() {
	os.Exit(_main())
}

func _main() int {
	const (
		help = "help"
	)

	command := help
	// options := []string{}

	if len(os.Args) > 1 {
		command = os.Args[1]
		// options = os.Args[2:]
	}

	switch command {
	case help:
		fmt.Print(usage)
	default:
		fmt.Print("Invalid Subcommands: Please Check help or typo.")
		return 1
	}
	return 0
}
