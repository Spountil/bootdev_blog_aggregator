package main

import (
	"bootdev_blog_aggregator/internal/config"
	"bootdev_blog_aggregator/internal/database"
	"bootdev_blog_aggregator/internal/rss"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type state struct {
	db   *database.Queries
	conf *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Error: the login handler expects one argument, the username")
	}

	ctx := context.Background()
	name := cmd.args[0]
	user, err := s.db.GetUser(ctx, name)
	if err != nil {
		return err
	}

	config.SetUser(*s.conf, user.Name)
	fmt.Printf("User %s has been logged!\n", user.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Error: the login handler expects one argument, the username")
	}

	ctx := context.Background()
	arg := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}
	user, err := s.db.CreateUser(ctx, arg)
	if err != nil {
		return err
	}

	config.SetUser(*s.conf, cmd.args[0])
	fmt.Printf("User %s has been set!\n", user.Name)
	fmt.Printf("User created with the following data; \n%v", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	err := s.db.DeleteUsers(ctx)
	if err != nil {
		return err
	}
	fmt.Println("Users table deleted successfully!")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	currentUser := s.conf.CurrentUserName
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range users {
		if user == currentUser {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}

func handlerAggregator(s *state, cmd command) error {
	feedURL := "https://www.wagslane.dev/index.xml"
	ctx := context.Background()
	rssFeed, err := rss.FetchFeed(ctx, feedURL)
	if err != nil {
		return err
	}
	fmt.Print(rssFeed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("Error: the addFeed handler expects two argument: the name and the URL.")
	}
	ctx := context.Background()
	userInfo, err := s.db.GetUser(ctx, s.conf.CurrentUserName)
	if err != nil {
		return err
	}
	userID := userInfo.ID

	arg := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    userID,
	}

	res, err := s.db.CreateFeed(ctx, arg)
	if err != nil {
		return err
	}

	fmt.Print(res)
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	fn, ok := c.cmds[cmd.name]
	if !ok {
		return fmt.Errorf("Function %s not found", cmd.name)
	}

	err := fn(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) error {
	if len(c.cmds) > 0 {
		_, ok := c.cmds[name]
		if ok {
			return fmt.Errorf("Function %s already registered", name)
		}
	}

	c.cmds[name] = f
	return nil
}
