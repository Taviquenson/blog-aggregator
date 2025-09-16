package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Taviquenson/gator/internal/config"
	"github.com/Taviquenson/gator/internal/database"
	_ "github.com/lib/pq"
) // Not used directly in the code, underscore
//   tells Go that it's imported it for its side
//   effects, not because it will be used explicitly.

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.Db_url)
	if err != nil {
		log.Fatalf("error opening connection to SQL database: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db) // type: *database.Queries

	programState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		CmdsMap: make(map[string]func(*state, command) error),
	}

	cmds.Register("login", handlerLogin)
	cmds.Register("register", handlerRegister)
	cmds.Register("reset", handlerReset)
	cmds.Register("users", handlerListUsers)
	cmds.Register("agg", handlerAgg)
	cmds.Register("addfeed", handlerAddfeed)
	cmds.Register("feeds", handlerFeeds)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmd := command{
		Name: args[1],
		Args: args[2:],
	}

	err = cmds.Run(programState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
