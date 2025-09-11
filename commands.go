package main

import (
	"errors"
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	CmdsMap map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("the login handler expects a single argument, the username")
	}
	username := cmd.Args[0] // <---- DOUBLE CHECK!
	err := s.cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Println("User has been set.")
	return nil
}

func (c *commands) Run(s *state, cmd command) error {
	myCmd, exists := c.CmdsMap[cmd.Name]
	if exists {
		err := myCmd(s, cmd)
		return err
	} else {
		return errors.New("can't run command, it doesn't exist")
	}
}

func (c *commands) Register(name string, f func(*state, command) error) {
	c.CmdsMap[name] = f
}
