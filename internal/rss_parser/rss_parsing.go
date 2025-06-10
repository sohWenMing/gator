package rssparser

import (
	"encoding/xml"
)

type RSSFeedData struct {
	Data ChannelData `xml:"channel"`
}

type ParsedRssFeed struct {
	Title       string
	Description string
	Language    string
	Link        string
}

type Link struct {
	Href     string `xml:"href,attr"`
	Rel      string `xml:"rel,attr"`
	Type     string `xml:"type,attr"`
	CharData string `xml:",chardata"`
}

type ChannelData struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Language    string `xml:"language"`
	Links       []Link `xml:"link"`
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
		Title:       rfd.Data.Title,
		Description: rfd.Data.Description,
		Language:    rfd.Data.Language,
		Link:        link,
	}
	return rSSFeed, nil

}
