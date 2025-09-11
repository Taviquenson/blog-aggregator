package config

import (
	"errors"
	"fmt"
)

type State struct {
	Cfg *Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	CmdsMap map[string]func(*State, Command) error
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("the login handler expects a single argument, the username")
	}
	username := cmd.Args[0] // <---- DOUBLE CHECK!
	err := s.Cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Println("User has been set.")
	return nil
}

func (c *Commands) Run(s *State, cmd Command) error {
	myCmd, exists := c.CmdsMap[cmd.Name]
	if exists {
		err := myCmd(s, cmd)
		return err
	} else {
		return errors.New("can't run command, it doesn't exist")
	}
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CmdsMap[name] = f
}
