package main

import (
	"log"
	"os"

	"github.com/sebasukodo/blog-aggregator/internal/config"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	cfgState := &state{
		cfg: cfg,
	}

	cmds := &commands{
		cmnds: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	osArgs := os.Args

	if len(osArgs) < 2 {
		log.Fatal("error, not enough arguments were provided")
	}

	cmdName := osArgs[1]

	cmdArguments := osArgs[2:]

	cmd := command{
		name:      cmdName,
		arguments: cmdArguments,
	}

	if err := cmds.run(cfgState, cmd); err != nil {
		log.Fatalf("error: %v", err)
	}

}
