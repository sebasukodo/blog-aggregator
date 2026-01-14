package main

import (
	"fmt"

	"github.com/sebasukodo/blog-aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	cmnds map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("no arguments were given")
	}
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("too many arguments were given")
	}

	if err := s.cfg.SetUser(cmd.arguments[0]); err != nil {
		return err
	}

	msg := fmt.Sprintf("User has been set to: %v", cmd.arguments[0])
	fmt.Println(msg)

	return nil

}

func (c *commands) run(s *state, cmd command) error {
	h, ok := c.cmnds[cmd.name]
	if !ok {
		return fmt.Errorf("command not found")
	}

	return h(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmnds[name] = f
}
