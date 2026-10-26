package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/DrowsyMissile/gator/internal/database"
	"github.com/google/uuid"
)

func handleFollow(s *State, cmd Command, user database.User) error {
	if err := requireArgs(cmd, "<feed_url>"); err != nil {
		return err
	}

	feed, err := s.DB.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}

	if err := followFeed(s, user.ID, feed.ID); err != nil {
		return err
	}

	fmt.Println("Feed follow created:")
	printFeedFollow(user.Name, feed.Name)
	return nil
}

func handleFollowing(s *State, cmd Command, user database.User) error {
	feedNames, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get feed follows: %w", err)
	}

	if len(feedNames) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	fmt.Printf("Feed follows for user %s:\n", user.Name)
	for _, name := range feedNames {
		fmt.Printf("* %s\n", name)
	}
	return nil
}

func handleUnfollow(s *State, cmd Command, user database.User) error {
	if err := requireArgs(cmd, "<feed_url>"); err != nil {
		return err
	}

	feed, err := s.DB.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	return s.DB.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
}

func followFeed(s *State, userID, feedID uuid.UUID) error {
	now := time.Now().UTC()
	err := s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    userID,
		FeedID:    feedID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}
	return nil
}

func printFeedFollow(userName, feedName string) {
	fmt.Printf("* User:          %s\n", userName)
	fmt.Printf("* Feed:          %s\n", feedName)
}
