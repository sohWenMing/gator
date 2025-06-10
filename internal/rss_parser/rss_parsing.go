package rssparser

import (
	"encoding/xml"
	"html"
)

type RSSFeedData struct {
	Data ChannelData `xml:"channel"`
}

type ParsedRssFeed struct {
	Title       string
	Description string
	Language    string
	Link        string
	Items       []RSSItem
}

type Link struct {
	Href     string `xml:"href,attr"`
	Rel      string `xml:"rel,attr"`
	Type     string `xml:"type,attr"`
	CharData string `xml:",chardata"`
}

type ChannelData struct {
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	Links       []Link    `xml:"link"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func ParseRSSXML(data []byte) (rSSFeed ParsedRssFeed, err error) {
	var rfd RSSFeedData
	err = xml.Unmarshal(data, &rfd)
	if err != nil {
		return ParsedRssFeed{}, err
	}
	link := ""
	for _, checkedLink := range rfd.Data.Links {
		if checkedLink.CharData != "" {
			link = checkedLink.CharData
			break
		} else {
			continue
		}
	}
	rSSFeed = ParsedRssFeed{
		Title:       html.UnescapeString(rfd.Data.Title),
		Description: html.UnescapeString(rfd.Data.Description),
		Language:    rfd.Data.Language,
		Link:        link,
		Items:       make([]RSSItem, 0, len(rfd.Data.Items)),
	}

	for _, item := range rfd.Data.Items {
		rSSFeed.Items = append(rSSFeed.Items, UnescapeStringsforRssItem(item))
	}
	return rSSFeed, nil
}

func UnescapeStringsforRssItem(item RSSItem) RSSItem {
	return RSSItem{
		Title:       html.UnescapeString(item.Title),
		Link:        item.Link,
		Description: html.UnescapeString(item.Description),
		PubDate:     item.PubDate,
	}
}
