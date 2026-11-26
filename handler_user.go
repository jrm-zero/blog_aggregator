package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jrm-zero/blog_aggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("One Argument Expected: Username")
	}

	usr, err := s.dBase.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		return errors.New("User not registered! Please register user to login.")
	}

	if err := s.cfg.SetUser(usr.Name); err != nil {
		return err
	}

	fmt.Printf("%s has been set as User\n", usr.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("One Argument Expected: Username")
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	}

	usr, err := s.dBase.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}

	if err := s.cfg.SetUser(usr.Name); err != nil {
		return err
	}
	fmt.Printf("Name: %s\nID: %v\nCreated At: %v\nUpdated At: %v\n", usr.Name, usr.ID, usr.CreatedAt, usr.UpdatedAt)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.dBase.DeleteAllUsers(context.Background()); err != nil {
		return err
	}
	fmt.Println("Database has been successfully reset!")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	usrs, err := s.dBase.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, usr := range usrs {
		if usr.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", usr.Name)
			continue
		}

		fmt.Printf("* %s\n", usr.Name)
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	feedURL := "https://www.wagslane.dev/index.xml"

	RSSFeed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}

	fmt.Printf("Title of Channel: %s\nDescription of Channel: %s\n\n", RSSFeed.Channel.Title, RSSFeed.Channel.Description)

	for _, item := range RSSFeed.Channel.Item {
		fmt.Printf("Title of Item: %s\nDescription of Item: %s\n\n", item.Title, item.Description)
	}
	return nil
}
