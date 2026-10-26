package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/DrowsyMissile/gator/internal/commands"
	"github.com/DrowsyMissile/gator/internal/config"
	"github.com/DrowsyMissile/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()

	s := &commands.State{DB: database.New(db), Cfg: &cfg}
	if err := commands.Run(s, commands.Command{Name: os.Args[1], Args: os.Args[2:]}); err != nil {
		log.Fatal(err)
	}
}
