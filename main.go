package main

import (
	"fmt"
	"log"
	"os"
	"github.com/jrm-zero/blog_aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	c, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", c)

	s := state{
		cfg: &c,
	}

	h := map[string]func(*state, command) error{}
	cmds := commands{
		handlers: h,
	}

	cmdLogin := command{
		name: "login",
		arguments: []string{},
	}

	err = cmds.register("login", handlerLogin)
	if err != nil {
		log.Fatalf("error registering command: %v", err)
	}

	cmdLnArgs := os.Args

	if len(cmdLnArgs) < 2 {
		fmt.Println("To few arguments entered!")
		os.Exit(1)
	}
	cmdLogin.arguments[0] = cmdLnArgs[1]
	err = cmds.run(&s, cmdLogin)
	if err != nil {
		log.Fatalf("error running command: %v", err)
	}


	c, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config again: %+v\n", c)
}