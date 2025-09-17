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
	username := s.cfg.CurrentUserName
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
	fmt.Println("Feed was created in database:")
	printFeed(feed, currUser)

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currUser.ID,
		FeedID:    feed.ID,
	}
	feedFollowRow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("couldn't create feedFollow: %w", err)
	}
	fmt.Printf("User %s now follows feed %s\n", feedFollowRow.UserName, feedFollowRow.FeedName)

	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %v", cmd.Name)
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Printf("Found %d feeds:\n", len(feeds))
	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get feed user: %w", err)
		}
		printFeed(feed, user)
	}
	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("ID: %v\nCreated_At: %v\nUpdated_At: %v\nName: %v\nURL: %v\nUser: %v\n\n",
		feed.ID, feed.CreatedAt, feed.UpdatedAt, feed.Name, feed.Url, user.Name)
}
