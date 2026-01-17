package main

import (
	"context"
	"fmt"
)

// This function is only for testing purposes and will be removed later
func handlerReset(s *state, cmd command) error {

	if err := s.db.ResetDatabase(context.Background()); err != nil {
		return fmt.Errorf("cannot reset database: %v", err)
	}

	fmt.Println("All exisiting users have been removed from database...")

	return nil

}
