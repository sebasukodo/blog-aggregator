package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sebasukodo/blog-aggregator/internal/config"
	"github.com/sebasukodo/blog-aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	cmnds map[string]func(*state, command) error
}

func handlerListUsers(s *state, cmd command) error {

	if len(cmd.arguments) != 0 {
		return fmt.Errorf("no arguments needed for this command")
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	var msg string
	for i := 0; i < len(users); i++ {
		if i > 0 {
			msg += "\n"
		}

		if users[i].Name == s.cfg.CurrentUserName {
			msg += fmt.Sprintf("* %v (current)", users[i].Name)
		} else {
			msg += fmt.Sprintf("* %v", users[i].Name)
		}
	}

	fmt.Println(msg)

	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("not enough arguments to login")
	}
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("too many arguments to login")
	}

	_, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("user was not found in database: %v", err)
	}

	if err := s.cfg.SetUser(cmd.arguments[0]); err != nil {
		return err
	}

	msg := fmt.Sprintf("User has been set to: %v", cmd.arguments[0])
	fmt.Println(msg)

	return nil

}

func handlerRegister(s *state, cmd command) error {

	if len(cmd.arguments) == 0 {
		return fmt.Errorf("not enough arguments to register a new user were given")
	}
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("too many arguments to register a new user were given")
	}

	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	}

	userData, err := s.db.CreateUser(context.Background(), newUser)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("User %v has been created", cmd.arguments[0])
	fmt.Println(msg)

	if err := handlerLogin(s, cmd); err != nil {
		return fmt.Errorf("could not create user %s: %w", cmd.arguments[0], err)
	}

	fmt.Println(userData)

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
