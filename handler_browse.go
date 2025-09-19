package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Taviquenson/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	postsLimit, err := 2, fmt.Errorf("")
	if len(cmd.Args) == 1 {
		postsLimit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("usage: %v [optional_number_of_posts]", cmd.Name)
		}
	} else if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %v [optional_number_of_posts]", cmd.Name)
	}

	postsForUserParams := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(postsLimit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), postsForUserParams)
	if err != nil {
		return fmt.Errorf("couldn't get posts for user %v: %w", user.Name, err)
	}
	err = printPosts(posts)
	if err != nil {
		return fmt.Errorf("error printing posts: %w", err)
	}
	return nil
}

func printPosts(posts []database.Post) error {
	if len(posts) > 0 {
		fmt.Println("These are your posts:")
	} else {
		fmt.Println("No posts were found. Make sure to aggregated your feeds")
	}
	for _, post := range posts {
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("Description: %s\n", post.Description)
		fmt.Printf("Published at: %v\n", post.PublishedAt)
	}
	return nil
}
