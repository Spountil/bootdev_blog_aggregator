package main

import (
	"bootdev_blog_aggregator/internal/config"
	"bootdev_blog_aggregator/internal/database"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	config, err := config.Read()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	db, err := sql.Open("postgres", config.Connection)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	s := state{
		db:   dbQueries,
		conf: &config,
	}

	cmds := commands{
		cmds: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)

	args := os.Args
	var cmdArgs []string
	if len(args) < 2 {
		fmt.Println("Error: Not enough arguments, expected 2.")
		os.Exit(1)
	} else if len(args) == 2 {
		cmdArgs = nil
	} else {
		cmdArgs = args[2:]
	}

	cmd := command{
		name: args[1],
		args: cmdArgs,
	}

	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
