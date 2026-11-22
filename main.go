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

	cmds.register("login", handlerLogin)

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
