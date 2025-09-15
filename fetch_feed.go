package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var body io.Reader = nil // For GET requests, the body istypically nil
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error creating request: %v", err)
	}
	// common practice to identify your program to the server
	req.Header.Add("User-Agent", "gator")
	// Send the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	// Process the response
	fmt.Printf("Response Status: %s\n", res.Status)
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error reading response body: %v", err)
	}
	// fmt.Printf("Response Body: %s\n", &resBody)
	var rssFeed RSSFeed
	if err := xml.Unmarshal(resBody, &rssFeed); err != nil {
		return &RSSFeed{}, fmt.Errorf("error unmarshalling XML: %v", err)
	}

	return &rssFeed, nil
}
