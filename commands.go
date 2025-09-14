package main

import (
	"errors"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	CmdsMap map[string]func(*state, command) error
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
