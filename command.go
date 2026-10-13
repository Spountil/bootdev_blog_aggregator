package main

import (
	"bootdev_blog_aggregator/internal/config"
	"bootdev_blog_aggregator/internal/database"
	"bootdev_blog_aggregator/internal/rss"
	"context"
	"database/sql"
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
	if len(cmd.args) < 1 {
		return fmt.Errorf("Error: the agg handler expects one argument: the time between requests.")
	}

	time_between_reqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("Error: the addFeed handler expects two argument: the name and the URL.")
	}

	arg := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	ctx := context.Background()
	_, err := s.db.CreateFeed(ctx, arg)
	if err != nil {
		fmt.Print(err)
		fmt.Print("Adding the information in the feed db")
	}

	cmd.args[0] = cmd.args[1]
	err = handlerFollow(s, cmd, user)
	if err != nil {
		return err
	}
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		userName, err := s.db.GetUserName(ctx, feed.UserID)
		if err != nil {
			return err
		}
		fmt.Print("\n")
		fmt.Printf("Name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("UserName: %s\n", userName)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Error: the follow handler expects one argument: the URL.")
	}

	ctx := context.Background()
	feedID, err := s.db.GetFeedsIDFromURL(ctx, cmd.args[0])
	if err != nil {
		return err
	}

	arg := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedID,
	}

	res, err := s.db.CreateFeedFollow(ctx, arg)
	if err != nil {
		return err
	}

	fmt.Printf("Name of the feed: %s\nCurrent User: %s\n", res.FeedName, s.conf.CurrentUserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	feedNames, err := s.db.GetFeedFollowsForUser(ctx, user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User %s follows:\n", s.conf.CurrentUserName)

	for _, feedName := range feedNames {
		fmt.Printf(" - %s\n", feedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Error: the unfollow handler expects one argument: the URL.")
	}

	ctx := context.Background()
	feedID, err := s.db.GetFeedsIDFromURL(ctx, cmd.args[0])
	if err != nil {
		return err
	}

	args := database.DeleteFeedFollowsParams{
		UserID: user.ID,
		FeedID: feedID,
	}

	err = s.db.DeleteFeedFollows(ctx, args)
	if err != nil {
		return err
	}

	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		user, err := s.db.GetUser(ctx, s.conf.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	nextFeed, err := s.db.GetNextFeedFetch(ctx)
	if err != nil {
		return err
	}

	args := database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ID: nextFeed.ID,
	}

	err = s.db.MarkFeedFetched(ctx, args)
	if err != nil {
		return nil
	}

	rssFeed, err := rss.FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return err
	}

	for _, item := range rssFeed.Channel.Item {
		fmt.Printf(" - %s\n", item.Title)
	}
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
