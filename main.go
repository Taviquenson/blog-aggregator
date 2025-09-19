package main

import (
	"context"
	"database/sql"
	"fmt"
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

var programState = &state{}
var cmd = command{}

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

	// programState := &state{
	// 	db:  dbQueries,
	// 	cfg: &cfg,
	// }
	programState.db = dbQueries
	programState.cfg = &cfg

	cmds := commands{
		CmdsMap: make(map[string]func(*state, command) error),
	}

	cmds.Register("login", handlerLogin)
	cmds.Register("register", handlerRegister)
	cmds.Register("reset", handlerReset)
	cmds.Register("users", handlerListUsers)
	cmds.Register("agg", handlerAgg)
	cmds.Register("addfeed", middlewareLoggedIn(handlerAddfeed))
	cmds.Register("feeds", handlerListFeeds)
	cmds.Register("follow", middlewareLoggedIn(handlerFollow))
	cmds.Register("following", middlewareLoggedIn(handlerListFeedFollows))
	cmds.Register("unfollow", middlewareLoggedIn(handlerUnfollow))

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	// cmd := command{
	// 	Name: args[1],
	// 	Args: args[2:],
	// }
	cmd.Name = args[1]
	cmd.Args = args[2:]

	err = cmds.Run(programState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}

// Takes a handler of the "logged in" type and returns a "normal" handler that can be registered
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(s *state, cmd command) error {
	return func(s *state, cmd command) error {
		username := programState.cfg.CurrentUserName
		user, err := programState.db.GetUser(context.Background(), username)
		if err != nil {
			return fmt.Errorf("no user logged in for function that requires one: %w", err)
		}
		err = handler(s, cmd, user)
		if err != nil {
			return err
		}
		return nil
	}
}
