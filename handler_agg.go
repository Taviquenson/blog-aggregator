package main

import (
	"context"
	"fmt"
	"html"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <time_between_reqs>\nValid time units are 's', 'm', 'h'", cmd.Name)
	}
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel() // Release resources associated with the context
	// channelLink := "https://www.wagslane.dev/index.xml"
	// rssFeed, err := fetchFeed(ctx, channelLink)
	// if err != nil {
	// 	log.Fatal(err) //  <------ CHANGE!!!!!!!!!!!****!!!!
	// }

	// unscp := html.UnescapeString
	// fmt.Printf("---RSSFeed---\nChannel Title: %s\nChannel Link: %v\nChannel Description: %s\n\nChannel Items:\n", unscp(rssFeed.Channel.Title), rssFeed.Channel.Link, unscp(rssFeed.Channel.Description))
	// for _, item := range rssFeed.Channel.Items {
	// 	fmt.Printf("-Title: %s\n-Link: %s\n-Description: %s\n-Publication Date: %s\n\n", unscp(item.Title), item.Link, unscp(item.Description), item.PubDate)
	// }
	// fmt.Printf("Feed: %+v\n", rssFeed) // Can use %+v to print a struct with its field names along with their corresponding values
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
	fmt.Printf("---RSSFeed: %s---\n", unscp(rssFeed.Channel.Title))
	for _, item := range rssFeed.Channel.Items {
		fmt.Printf("- Post Title: %s\n", unscp(item.Title))
	}

	return nil
}
