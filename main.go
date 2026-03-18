package main

import (
	"fmt"
	"log"
	"os"

	"github.com/CamilleOnoda/blog-aggregator/internal/config"
)

type state struct {
	*config.Config
}

type command struct {
	name string
	args []string
}

type CLIcommands struct {
	cmd map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Expected to receive a username after 'login' command")
	}

	username := cmd.args[0]
	err := config.SetUser(username)
	if err != nil {
		return fmt.Errorf("Error setting user: %v", err)
	}

	s.Current_user_name = username
	fmt.Println("User set to:", username)
	return nil
}

func (c *CLIcommands) run(s *state, cmd command) error {
	if handler, exists := c.cmd[cmd.name]; exists {
		return handler(s, cmd)
	} else {
		return fmt.Errorf("Unknown command: %s", cmd.name)
	}
}

// This function registers a command and its handler to the CLIcommands struct
func (c *CLIcommands) register(name string, f func(*state, command) error) {
	c.cmd[name] = f
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("Error reading config:", err)
	}

	s := &state{&cfg}
	cliCommands := &CLIcommands{cmd: make(map[string]func(*state, command) error)}
	cliCommands.register("login", handlerLogin)

	cliArgs := os.Args[1:]
	if len(cliArgs) < 2 {
		log.Fatal(fmt.Errorf(
			"Expected at least 2 arguments: command name and its arguments"))
	}

	cmd := command{
		name: cliArgs[0],
		args: cliArgs[1:],
	}
	if err := cliCommands.run(s, cmd); err != nil {
		log.Fatal(err)
	}
}
