package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/DrowsyMissile/gator/internal/config"
	"github.com/DrowsyMissile/gator/internal/database"
)

type State struct {
	DB  *database.Queries
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
}

type handler func(*State, Command) error

var registry = map[string]handler{
	"register":  handleRegister,
	"login":     handleLogin,
	"reset":     handleReset,
	"users":     handleUsers,
	"agg":       handleAgg,
	"addfeed":   loggedIn(handleAddFeed),
	"feeds":     handleFeeds,
	"follow":    loggedIn(handleFollow),
	"following": loggedIn(handleFollowing),
	"unfollow":  loggedIn(handleUnfollow),
}

func Run(s *State, cmd Command) error {
	f, ok := registry[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}

func requireArgs(cmd Command, usage ...string) error {
	if len(cmd.Args) != len(usage) {
		return fmt.Errorf("usage: %s %s", cmd.Name, strings.Join(usage, " "))
	}
	return nil
}
