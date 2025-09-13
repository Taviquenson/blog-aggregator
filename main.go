package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Taviquenson/gator/internal/config"
	"github.com/Taviquenson/gator/internal/database"
	_ "github.com/lib/pq"
) // not used directly in the code
// Underscore tells Go that it's imported it for its
// side effects, not because it will be used explicitly.

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	programState := &state{
		cfg: &cfg,
	}

	db, err := sql.Open("postgres", programState.cfg.Db_url)
	if err != nil {
		log.Fatalf("error opening connection to SQL database: %v", err)
	}
	dbQueries := database.New(db) // type: *database.Queries
	programState.db = dbQueries

	var cmds = commands{
		CmdsMap: make(map[string]func(*state, command) error),
	}

	cmds.Register("login", handlerLogin)
	cmds.Register("register", handlerRegister)

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Not enough arguments provided.")
	}

	cmd := command{
		Name: args[1],
		Args: args[2:],
	}

	err = cmds.Run(programState, cmd)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("Read config again: %+v\n", cfg)
}
