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

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.arguments) < 2 {
		return errors.New("Two Arguments Expected: go run . addfeed [Feed Name] [Feed URL]")
	}

	currentUsr, err := s.dBase.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		newErr := fmt.Errorf("Error Getting User: %v\n", err)
		return newErr
	}

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    currentUsr.ID,
	}

	feed, err := s.dBase.CreateFeed(context.Background(), feedParams)
	if err != nil {
		newErr := fmt.Errorf("Error Creating Feed: %v\n", err)
		return newErr
	}

	feedID, err := s.dBase.GetFeedIDByURL(context.Background(), cmd.arguments[1])
	if err != nil {
		newErr := fmt.Errorf("Error Getting Feed's ID: %v\n", err)
		return newErr
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUsr.ID,
		FeedID:    feedID,
	}

	feedFollow, err := s.dBase.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		newErr := fmt.Errorf("Error Creating Feed Follow While Creating New Feed: %v\n", err)
		return newErr
	}

	fmt.Printf("Feed Name: %s\nFeed URL: %s\nCreated At: %v\nUpdated At: %v\nUser ID: %s\n", feed.Name, feed.Url, feed.CreatedAt, feed.UpdatedAt, feed.UserID)
	fmt.Printf("%s is also now following %s", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.dBase.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		usr, err := s.dBase.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return errors.New("User not found when retreiving feeds.")
		}
		fmt.Printf("* Name: %s, URL: %s, User: %s\n", feed.Name, feed.Url, usr)
	}
	return nil
}

func handlerFollowFeed(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("One Argument Expected: Feed URL")
	}

	currentUsr, err := s.dBase.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	feedID, err := s.dBase.GetFeedIDByURL(context.Background(), cmd.arguments[0])
	if err != nil {
		return err
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUsr.ID,
		FeedID:    feedID,
	}

	feedFollow, err := s.dBase.CreateFeedFollow(context.Background(), feedFollowParams)

	fmt.Printf("ID of Row: %v\nCreated At: %v\nUpdated At: %v\nWho is Following: %s\nName of Feed: %s\n", feedFollow.ID, feedFollow.CreatedAt, feedFollow.UpdatedAt, feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.arguments) > 0 {
		return errors.New("No arguments expected for following")
	}

	currentUsr, err := s.dBase.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	feedFollows, err := s.dBase.GetFeedFollowsForUser(context.Background(), currentUsr.ID)
	if err != nil {
		return err
	}

	fmt.Printf("%s follows these feeds:\n", currentUsr.Name)
	for _, followedFeed := range feedFollows {
		fmt.Printf("* %s\n", followedFeed.FeedName)
	}

	return nil
}
