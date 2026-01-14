package main

import (
	"fmt"
	"log"

	"github.com/sebasukodo/blog-aggregator/internal/config"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	if err = cfg.SetUser("sebasukodo"); err != nil {
		log.Fatalf("error writing config file: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	fmt.Println("DB URL:", cfg.DBURL)
	fmt.Println("Current user:", cfg.CurrentUserName)

}
