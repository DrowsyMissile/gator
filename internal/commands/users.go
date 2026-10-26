package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/DrowsyMissile/gator/internal/database"
	"github.com/google/uuid"
)

func handleRegister(s *State, cmd Command) error {
	if err := requireArgs(cmd, "<name>"); err != nil {
		return err
	}

	now := time.Now().UTC()
	user, err := s.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.Args[0],
	})
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	if err := s.Cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User created successfully:")
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
	return nil
}

func handleLogin(s *State, cmd Command) error {
	if err := requireArgs(cmd, "<name>"); err != nil {
		return err
	}
	name := cmd.Args[0]

	if _, err := s.DB.GetUser(context.Background(), name); err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	if err := s.Cfg.SetUser(name); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User switched successfully!")
	return nil
}

func handleUsers(s *State, cmd Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't list users: %w", err)
	}
	for _, user := range users {
		suffix := ""
		if user.Name == s.Cfg.CurrentUserName {
			suffix = " (current)"
		}
		fmt.Printf("* %s%s\n", user.Name, suffix)
	}
	return nil
}

func handleReset(s *State, cmd Command) error {
	if err := s.DB.DeleteUsers(context.Background()); err != nil {
		return fmt.Errorf("couldn't delete users: %w", err)
	}
	fmt.Println("Database reset successfully!")
	return nil
}
