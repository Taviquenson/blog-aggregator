package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Taviquenson/gator/internal/database"
	"github.com/google/uuid"
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
	username := cmd.Args[0]
	user, err := s.db.GetUser(context.Background(), username)
	if user == (database.User{}) {
		fmt.Printf("Username '%s' doesn't exist in the database.\n", username)
		os.Exit(1)
	}
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("User has been set to %s\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("the register handler expects a single argument, the user's name")
	}
	name := cmd.Args[0]
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}
	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("User was created in database:\nID: %v\nCreatedAt: %v\nUpdatedAt: %v\nName: %v\n", user.ID, user.CreatedAt, user.UpdatedAt, user.Name)
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Printf("User has been set to %s.\n", user.Name)

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
