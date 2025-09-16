package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Taviquenson/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddfeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %v <name> <url>", cmd.Name)
	}
	feedName, url := cmd.Args[0], cmd.Args[1]
	username := s.cfg.Current_user_name
	currUser, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       url,
		UserID:    currUser.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}
	printFeed(feed)
	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("Feed was created in database:\nID: %v\nCreatedAt: %v\nUpdatedAt: %v\nName: %v\nURL: %v\nUser_ID: %v\n",
		feed.ID, feed.CreatedAt, feed.UpdatedAt, feed.Name, feed.Url, feed.UserID)
}
