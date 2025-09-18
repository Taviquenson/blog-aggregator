package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Taviquenson/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <url>", cmd.Name)
	}
	url := cmd.Args[0]
	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't find feed: %w", err)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollowRow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("couldn't create feedFollow: %w", err)
	}
	printFeedFollow(feedFollowRow.UserName, feedFollowRow.FeedName)

	return nil
}

func handlerListFeedFollows(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %v", cmd.Name)
	}
	feedFollowsForUserRows, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("couldn't find feeds for user %s: %w", user.Name, err)
	}
	fmt.Printf("Feeds followed by %s:\n", user.Name)
	for _, feedFollow := range feedFollowsForUserRows {
		fmt.Printf("- %s\n", feedFollow.FeedName)
	}
	return nil
}

func printFeedFollow(userName, feedName string) {
	fmt.Printf("User %s now follows feed %s\n", userName, feedName)
}
