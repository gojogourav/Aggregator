package main

import (
	"context"
	"database/sql"
	db "github/gojogourav/RSSAggregator/db/sqlc"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

func startScraping(db *db.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on  %v goroutines every %s duration", concurrency, timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)

		if err != nil {
			log.Println("Error fetching feeds : ", err)
			return
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(database *db.Queries, wg *sync.WaitGroup, feed db.Feed) {
	defer wg.Done()

	RSSFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed :", err)
		return
	}

	for _, item := range RSSFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		pubAt := sql.NullTime{}
		if item.PubDate != "" {
			t, err := time.Parse(time.RFC1123Z, item.PubDate)
			if err != nil {
				log.Printf("Failed to parse time %v", err)
				continue
			}

			pubAt.Time = t
			pubAt.Valid = true

		}

		_, err := database.CreatePost(context.Background(), db.CreatePostParams{
			ID:          uuid.New(),
			UpdatedAt:   time.Now(),
			CreatedAt:   time.Now(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("Failed to create post", err)

		}
	}

	log.Printf("Feed %s collected, %v posts found ", feed.Name, len(RSSFeed.Channel.Item))

	_, err = database.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched : ", err)
		return
	}
}
