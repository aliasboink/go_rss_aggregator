package main

import (
	"context"
	"encoding/xml"
	"log"
	"net/http"
	"rss/internal/database"
	"sync"
	"time"

	"github.com/google/uuid"
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
	PublishedAt string `xml:"pubDate"`
}

func fetchRSS(rssUrl string) (RSSFeed, error) {
	rss, err := http.Get(rssUrl)
	if err != nil {
		log.Println(err)
		return RSSFeed{}, err
	}
	decoder := xml.NewDecoder(rss.Body)
	rssFeed := RSSFeed{}
	err = decoder.Decode(&rssFeed)
	if err != nil {
		log.Printf("[fetchRSS] Error with %v: %v", rssUrl, err)
		return RSSFeed{}, err
	}
	return rssFeed, nil
}

func rssThiefWorker(db *database.Queries, interval time.Duration, numberOfFeeds int32) {
	log.Println("[WORKER] Worker started...")
	for range time.Tick(time.Second * time.Duration(interval)) {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), numberOfFeeds)
		if err != nil {
			log.Println("[WORKER] Error fetching feeds: ", err)
			continue
		}
		count := 0
		var wg sync.WaitGroup
		for _, feed := range feeds {
			wg.Add(1)
			log.Println("[WORKER] Work started on " + feed.Name)
			go func(feed database.Feed) {
				defer wg.Done()
				rssFeed, err := fetchRSS(feed.Url)
				if err != nil {
					log.Printf("[WORKER] Error fetching rss with url %v: %v", feed.Url, err)
					return
				}
				log.Println("[WORKER] URL of RSS: " + feed.Url)
				for _, item := range rssFeed.Channel.Items {
					// log.Println("[WORKER] Link of RSS item: " + item.Link)
					count++

					postParams, err := postToDatabasePost(Post{
						ID:          uuid.New(),
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
						Title:       item.Title,
						Url:         item.Link,
						Description: item.Description,
						PublishedAt: item.PublishedAt,
						FeedID:      feed.ID,
					})
					if err != nil {
						log.Println("[WORKER] Error parsing time: ", err)
						continue
					}

					_, err = db.CreatePost(context.Background(), postParams)
					if err != nil {
						if err.Error() != "pq: duplicate key value violates unique constraint \"posts_url_key\"" {
							log.Println("[WORKER] Error creating post: ", err.Error())
						}
						continue
					}
				}
			}(feed)
		}
		wg.Wait()
		log.Printf("[WORKER] Feeds done fetching %d feeds...", count)
	}
}
