package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/DrowsyMissile/gator/internal/database"
	"github.com/DrowsyMissile/gator/internal/rss"
	"github.com/google/uuid"
)

const separator = "====================================="

func handleAgg(s *State, cmd Command) error {
	feed, err := rss.Fetch(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}
	fmt.Printf("Feed: %+v\n", feed)
	return nil
}

func handleAddFeed(s *State, cmd Command, user database.User) error {
	if err := requireArgs(cmd, "<name>", "<url>"); err != nil {
		return err
	}

	now := time.Now().UTC()
	feed, err := s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	if err := followFeed(s, user.ID, feed.ID); err != nil {
		return err
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed, user.Name)
	fmt.Println()
	fmt.Println("Feed followed successfully:")
	printFeedFollow(user.Name, feed.Name)
	fmt.Println(separator)
	return nil
}

func handleFeeds(s *State, cmd Command) error {
	rows, err := s.DB.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}

	if len(rows) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Printf("Found %d feeds:\n", len(rows))
	for _, row := range rows {
		printFeed(row.Feed, row.UserName)
		fmt.Println(separator)
	}
	return nil
}

func printFeed(feed database.Feed, userName string) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", userName)
}
