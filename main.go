package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/sebasukodo/blog-aggregator/internal/config"
	"github.com/sebasukodo/blog-aggregator/internal/database"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	cfgState := &state{
		cfg: cfg,
	}

	db, err := sql.Open("postgres", cfgState.cfg.DBURL)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	dbQueries := database.New(db)

	cfgState.db = dbQueries

	cmds := &commands{
		cmnds: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)

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
		fmt.Println(err)
		os.Exit(1)
	}

}
