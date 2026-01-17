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
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerFetch)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerListAllFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))

	if len(os.Args) < 2 {
		log.Fatal("error, not enough arguments were provided")
	}

	cmd := command{
		name:      os.Args[1],
		arguments: os.Args[2:],
	}

	if err := cmds.run(cfgState, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
