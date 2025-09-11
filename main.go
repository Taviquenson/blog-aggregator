package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Taviquenson/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	// err = cfg.SetUser("tavo")
	// if err != nil {
	// 	log.Fatalf("couldn't set current user: %v", err)
	// }

	// cfg, err = config.Read()
	// if err != nil {
	// 	log.Fatalf("error reading config: %v", err)
	// }

	programState := &state{
		cfg: &cfg,
	}
	// fmt.Printf("Stored state: %v\n", state)

	var cmds = commands{
		CmdsMap: make(map[string]func(*state, command) error),
	}

	cmds.Register("login", handlerLogin)

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
