package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/KamogeloMahlake/go_blog_aggregator/internal/config"
	"github.com/KamogeloMahlake/go_blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)
	
	db, err := sql.Open("postgres", cfg.DBURL)
	dbQueries := database.New(db)

	s := &state{
		cfg: &cfg,
		db: dbQueries,
	}
	cmds := commands{
		handlerFuctions: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handleRegister)
	if len(os.Args) < 2 {
		log.Fatalf("No arguments provided")
		return
	}

	err = cmds.run(s, command{name: os.Args[1], arguments: os.Args[2:]})

	if err != nil {
		log.Fatalf("Error: %v", err)
	}


	fmt.Printf("Read config again: %+v\n", cfg)

}
