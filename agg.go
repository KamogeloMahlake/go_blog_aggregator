package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/KamogeloMahlake/go_blog_aggregator/internal/database"
	"github.com/google/uuid"
)



func handleAgg(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("Missing/Incorrect argument")
	}
	t, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		return err
	}

	var rssFeed RSSFeed
	for ticker := time.NewTicker(t); ; <-ticker.C {
		if err = scrapeFeed(s, &rssFeed); err != nil {
			return err
		}
	}
}


func scrapeFeed(s *state, rssFeed *RSSFeed) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	mark, err := s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID: feed.ID,
		UpdatedAt: time.Now(),
	})

	if mark.UpdatedAt.Compare(mark.LastFetchedAt.Time) != 0 {
		return fmt.Errorf("Feed not updated")
	}
	rssFeed, err = fetchFeed(context.Background(), mark.Url)
	if err != nil {
		return err
	}

	for _, item := range rssFeed.Channel.Item{
		description := sql.NullString{
			String: item.Description,
			Valid: true,
		}

		if description.String == "" {
			description.Valid = false
		}

		pub, err := time.Parse("Fri, 01 Jan 1970 00:00:00", item.PubDate)
		if err != nil {
			return err
		}
	  _, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
		  ID: uuid.New(), 
		  CreatedAt: time.Now(),
		  Title: item.Title,
		  Url: item.Link,
		  Description: description.String,
		  PublishedAt: pub,
		  FeedID: feed.ID,
	  })
		
		if err != nil {
			return err
		}
	}
	return nil
}
