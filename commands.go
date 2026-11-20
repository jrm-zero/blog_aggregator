package main

import (
	"errors"
	"fmt"
)

type command struct {
	name string
	arguments []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if s.cfg == nil {
		return errors.New("Config file does not exist")
	}

	handler, ok := c.handlers[cmd.name]
	if ok != true {
		return errors.New("Command does not exist")
	}

	handler(s, cmd)
	return nil

}

func (c *commands) register(name string, f func(*state, command) error) error {
	_, ok := c.handlers[name]
	if ok != false {
		err := fmt.Errorf("%s exists in commands struct", name)
		return err
	}

	c.handlers[name] = f
	return nil
}