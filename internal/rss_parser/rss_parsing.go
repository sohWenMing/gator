package rssparser

import (
	"encoding/xml"
	"fmt"
	"html"
	"strings"
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

func (f *ParsedRssFeed) String() string {
	var b strings.Builder
	delimiter := fmt.Sprintln(fmt.Sprintf("+%s+", strings.Repeat("-", 45)))
	itemDelimiter := fmt.Sprintln(fmt.Sprintf("+%s+", strings.Repeat("-", 30)))
	titleString := fmt.Sprintln("Title: ", f.Title)
	DescriptionString := fmt.Sprintln("Description", f.Description)
	languageString := fmt.Sprintln("Language", f.Language)
	linkString := fmt.Sprintln("Link", f.Link)
	blankString := fmt.Sprintln("")

	b.WriteString(delimiter)
	b.WriteString(blankString)
	b.WriteString(titleString)
	b.WriteString(DescriptionString)
	b.WriteString(languageString)
	b.WriteString(linkString)

	for _, item := range f.Items {
		itemTitleString := fmt.Sprintln("News Item: ", item.Title)
		itemLinkString := fmt.Sprintln("Link: ", item.Link)
		itemDescriptionString := fmt.Sprintln("Description: ", item.Description)
		itemDateString := fmt.Sprintln("Date: ", item.PubDate)
		b.WriteString(itemDelimiter)
		b.WriteString(itemTitleString)
		b.WriteString(blankString)
		b.WriteString(itemLinkString)
		b.WriteString(blankString)
		b.WriteString(strings.ReplaceAll(strings.ReplaceAll(itemDescriptionString, "<p>", ""), "</p>", ""))
		b.WriteString(blankString)
		b.WriteString(itemDateString)
		b.WriteString(blankString)
		b.WriteString(blankString)
		b.WriteString(itemDelimiter)

	}
	b.WriteString(blankString)
	b.WriteString(delimiter)

	return b.String()

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
