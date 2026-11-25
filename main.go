package main

import (
	"fmt"
	"log"
	"os"

	"database/sql"

	"github.com/jrm-zero/blog_aggregator/internal/config"
	"github.com/jrm-zero/blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg   *config.Config
	dBase *database.Queries
}

func main() {
	c, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", c)

	dbURL := c.DBURL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening data base: %v", err)
	}

	s := state{
		cfg:   &c,
		dBase: database.New(db),
	}

	h := map[string]func(*state, command) error{}
	cmds := commands{
		handlers: h,
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)

	cmdLnArgs := os.Args

	if len(cmdLnArgs) < 2 {
		fmt.Println("Too few arguments entered!")
		os.Exit(1)
	}

	newCmd := command{
		name:      cmdLnArgs[1],
		arguments: cmdLnArgs[2:],
	}

	err = cmds.run(&s, newCmd)
	if err != nil {
		log.Fatalf("error running command: %v", err)
	}

	c, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config again: %+v\n", c)
}
