package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ScriptedPranav/rss-aggregator/internal/database"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s duration\n", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(),int32(concurrency))

		if err != nil {
			log.Printf("Error fetching feeds %v\n",err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error Marking feed as fetched")
		return
	}
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:",err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		log.Println("Found post",item.Title)
	}
	log.Printf("Feed %s collected, %v posts found",feed.Name, len(rssFeed.Channel.Item))
}