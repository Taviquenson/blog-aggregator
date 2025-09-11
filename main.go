package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Taviquenson/gator/internal/config"
)

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

	state := config.State{
		Cfg: &cfg,
	}
	// fmt.Printf("Stored state: %v\n", state)

	var cmds = config.Commands{
		CmdsMap: make(map[string]func(*config.State, config.Command) error),
	}

	cmds.Register("login", config.HandlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Not enough arguments provided.")
	}

	cmd := config.Command{
		Name: args[1],
		Args: args[2:],
	}

	err = cmds.Run(&state, cmd)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("Read config again: %+v\n", cfg)
}
