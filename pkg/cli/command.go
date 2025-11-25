package cli

import (
	"fmt"
)

// create a command interface for all commands to implement
type Command interface {
	Name() string                // command name
	Description() string         // command description
	Usage() string               // command usage
	Execute(args []string) error // execute the command with arguments
}

// create a BaseCommand struct to embed in all commands.
type BaseCommand struct {
	name        string
	description string
	usage       string
	handler     func(args []string) error
}

// create a command error struct
type CommandError struct {
	Command string
	Err     error
}

// create a newbasecommand method
func NewBaseCommand(name, description, usage string, handler func(args []string) error) *BaseCommand {
	return &BaseCommand{
		name:        name,
		description: description,
		usage:       usage,
		handler:     handler,
	}
}

// create name method
func (c *BaseCommand) Name() string {
	return c.name
}

// create description method
func (c *BaseCommand) Description() string {
	return c.description
}

// create usage method
func (c *BaseCommand) Usage() string {
	return c.usage
}

// create execute method
func (c *BaseCommand) Execute(args []string) error {
	return c.handler(args)
}

// create an error method for command error
func (e *CommandError) Error() string {
	return fmt.Sprintf("command '%s': %v", e.Command, e.Err)
}

// create a unwrap method for command error
func (e *CommandError) Unwrap() error {
	return e.Err
}
