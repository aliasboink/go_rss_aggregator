package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

// This is the first successful attempt and I don't like it all that much.
type RSSFeed struct {
	Channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Items       []Item `xml:"item"`
	} `xml:"channel"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
}

func fetchRSS(rssUrl string) error {
	rss, err := http.Get(rssUrl)
	if err != nil {
		log.Println(err)
		return err
	}
	decoder := xml.NewDecoder(rss.Body)
	rssFeed := RSSFeed{}
	err = decoder.Decode(&rssFeed)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(rssUrl)
	log.Println(err)
	for _, item := range rssFeed.Channel.Items {
		fmt.Println(item)
	}
	return nil
}
