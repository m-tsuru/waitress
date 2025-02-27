package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/emiago/sipgo/sip"
	"github.com/m-tsuru/waitress/lib"
)

//go:embed VERSION
var version string

var usage = fmt.Sprintf(`
waitress %s
`, version)[1:]

func main() {
	sip.SIPDebug = true
	os.Exit(_main())
}

func _main() int {
	const (
		help = "help"
		register = "register"
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
	case register:
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()
		conf, regOpts, err := lib.ParseClientAccount("config.ini")
		if err != nil {
			log.Fatalf("Configuration Parse Error: %v\n", err)
		}
		forwIP, err := lib.ParseForwardIP("config.ini")
		err = lib.Setup(ctx, *conf, *forwIP, *regOpts)
		if err != nil {
			log.Fatalf("Configuration Parse Error: %v\n", err)
		}
	default:
		fmt.Print("Invalid Subcommands: Please Check help or typo.")
		return 1
	}
	return 0
}
