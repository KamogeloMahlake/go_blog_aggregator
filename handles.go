package main

import (
	"context"
	"fmt"
	"time"

	"github.com/KamogeloMahlake/go_blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("Username must be provided")
	}

	if len(cmd.arguments) != 1 {
		return fmt.Errorf("login expects a single argument, the username")
	}

	_, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}

	fmt.Printf("%s has been set as user.\n", cmd.arguments[0])
	return nil

}

func handleRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("Username must be provided")
	}

	if len(cmd.arguments) != 1 {
		return  fmt.Errorf("login expects a single argument, the username")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.arguments[0]})
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Printf("User Registered: %+v", user)
	return nil
}
