package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, _ command) error {
	// The underscore above indicates we are matching a
	// function signature for handlers but do not intend
	// to use that value of type command in this function
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't delete database rows: %w", err)
	}
	fmt.Println("Successfully reset the database.")
	return nil
}
