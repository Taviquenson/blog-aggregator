package main

import (
	"context"
	"fmt"
	"html"
	"log"
	"time"
)

func handlerAgg(s *state, _ command) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Release resources associated with the context
	channelLink := "https://www.wagslane.dev/index.xml"
	rssFeed, err := fetchFeed(ctx, channelLink)
	if err != nil {
		log.Fatal(err)
	}

	unscp := html.UnescapeString
	fmt.Printf("---RSSFeed---\nChannel Title: %s\nChannel Link: %v\nChannel Description: %s\n\nChannel Items:\n", unscp(rssFeed.Channel.Title), rssFeed.Channel.Link, unscp(rssFeed.Channel.Description))
	for _, item := range rssFeed.Channel.Items {
		fmt.Printf("-Title: %s\n-Link: %s\n-Description: %s\n-Publication Date: %s\n\n", unscp(item.Title), item.Link, unscp(item.Description), item.PubDate)
	}
	return nil
}
