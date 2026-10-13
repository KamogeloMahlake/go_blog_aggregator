package main

import "fmt"

type command struct {
	name string
	arguments []string
}

type commands struct {
	handlerFuctions map[string]func(*state, command) error
}


func (c *commands) run(s *state, cmd command) error {
	f, ok := c.handlerFuctions[cmd.name]

	if !ok {
		return fmt.Errorf("Command does not exist")
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlerFuctions[name] = f
}
