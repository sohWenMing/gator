package rssparser

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestRssParser(t *testing.T) {
	type test struct {
		name                string
		url                 string
		expectedTitle       string
		expectedDescription string
		expectedLanguage    string
		expectedLink        string
	}

	tests := []test{
		{
			"testing boot.dev rss",
			"https://www.wagslane.dev/index.xml",

			"Lane's Blog",
			"Recent content on Lane's Blog",
			"en-us",
			"https://wagslane.dev/",
		},
		{
			"testing straits times",
			"https://www.straitstimes.com/news/business/rss.xml",
			"The Straits Times Business News",
			"The Straits Times" +
				" - " + "Get exclusive stories, in-depth " +
				"analyses and award-winning multimedia content about " +
				"Singapore, Asia and the world.",
			"en",
			"https://www.straitstimes.com/",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, test.url, nil)
			if err != nil {
				t.Errorf("didn't expect error, got %v", err)
			}
			req.Header.Set("User-Agent", "gator")
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("didn't expect error, got %v", err)
			}
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("didn't expect error, got %v", err)
			}
			defer res.Body.Close()
			got, err := ParseRSSXML(data)
			if err != nil {
				t.Errorf("didn't expect error, got %v", err)
			}
			if got.Title != test.expectedTitle {
				t.Errorf("\ngot: %s\nwant: %s", got.Title, test.expectedTitle)
			}
			if got.Description != test.expectedDescription {
				t.Errorf("\ngot: %s\nwant: %s", got.Description, test.expectedDescription)
			}
			if got.Language != test.expectedLanguage {
				t.Errorf("\ngot: %s\nwant: %s", got.Language, test.expectedLanguage)
			}
			if got.Link != test.expectedLink {
				t.Errorf("\ngot: %s\nwant: %s", got.Link, test.expectedLink)
			}
			fmt.Println("")
			fmt.Println("Items: ", got.Items)
			fmt.Println("")
			fmt.Println("")
		})
	}
}
