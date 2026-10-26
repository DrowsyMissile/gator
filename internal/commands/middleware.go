package commands

import (
	"context"

	"github.com/DrowsyMissile/gator/internal/database"
)

type userHandler func(*State, Command, database.User) error

func loggedIn(h userHandler) handler {
	return func(s *State, cmd Command) error {
		user, err := s.DB.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return h(s, cmd, user)
	}
}
