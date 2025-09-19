package main

import (
	"context"
	"fmt"
	"html"
	"time"

	"github.com/Taviquenson/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/microcosm-cc/bluemonday"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <time_between_reqs>\nExamples of valid time units are 's', 'm', 'h'", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error parsing the duration of arg '%s'", cmd.Args[0])
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get the next feed to fetch: %w", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("couldn't mark feed as fetched: %w", err)
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("couldn't fetch the feed: %w", err)
	}
	unscp := html.UnescapeString
	p := bluemonday.StrictPolicy() // removes all tags
	for _, item := range rssFeed.Channel.Items {
		layout := "Mon, 02 Jan 2006 15:04:05 -0700"
		parsedTime, err := time.Parse(layout, item.PubDate)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("error parsing time: %w", err)
		}
		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       unscp(p.Sanitize(item.Title)),
			Url:         item.Link,
			Description: unscp(p.Sanitize(item.Description)),
			PublishedAt: parsedTime,
			FeedID:      feed.ID,
		}
		post, err := s.db.CreatePost(context.Background(), postParams)
		if err != nil {
			pqError, ok := err.(*pq.Error)
			if ok && pqError.Code == "23505" {
				continue
			}
			return fmt.Errorf("error saving post %s: %w", post.Title, err)
		}
	}

	return nil
}
