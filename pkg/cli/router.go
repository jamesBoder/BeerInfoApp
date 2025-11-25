package cli

import (
	"fmt"
	"strings"
)

// create a router struct to hold commands
type Router struct {
	commands      map[string]Command // map of command name to command
	aliases       map[string]string  // map of alias to command name
	multiWordCmds []string           // list of multi-word commands
	defaultCmd    Command            // default command to execute if no command is specified
}

// create a new router method
func NewRouter() *Router {
	return &Router{
		commands:      make(map[string]Command),
		aliases:       make(map[string]string),
		multiWordCmds: []string{},
	}
}

// create a register command method
func (r *Router) RegisterCommand(cmd Command) error {
	// validates cmd has non-empty name
	if cmd.Name() == "" {
		return fmt.Errorf("command name cannot be empty")
	}

	// return error if name exists
	if _, exists := r.commands[cmd.Name()]; exists {
		return fmt.Errorf("command already registered: %s", cmd.Name())
	}

	// register the command
	r.commands[cmd.Name()] = cmd

	return nil

}

// create a method that registers an alias for a command
func (r *Router) RegisterAlias(alias string, cmdName string) error {
	// creates an alias for a command
	if _, exists := r.commands[cmdName]; !exists {
		return fmt.Errorf("command not found: %s", cmdName)
	}
	r.aliases[alias] = cmdName
	return nil
}

// create a method to register multi-word commands
func (r *Router) RegisterMultiWordCommand(cmdName string) error {
	// marks a command as multi-word
	if _, exists := r.commands[cmdName]; !exists {
		return fmt.Errorf("command not found: %s", cmdName)
	}
	r.multiWordCmds = append(r.multiWordCmds, cmdName)
	return nil
}

// create a set default command method
func (r *Router) SetDefaultCommand(cmd Command) {
	r.defaultCmd = cmd
}

// create a Execute method
func (r *Router) Execute(input string) error {
	// parse input into command and args
	parts := strings.Fields(input)

	// handle empty input
	if len(parts) == 0 {
		if r.defaultCmd != nil {
			return r.defaultCmd.Execute([]string{})
		}
		return fmt.Errorf("no command provided")
	}

	// finds matching command
	var cmdName string
	var args []string

	// check for multi-word commands first
	for _, mwCmd := range r.multiWordCmds {
		mwParts := strings.Fields(mwCmd)
		if len(parts) >= len(mwParts) && strings.Join(parts[:len(mwParts)], " ") == mwCmd {
			cmdName = mwCmd
			args = parts[len(mwParts):]
			break
		}
	}

	// if no multi-word command matched, use first part as command name
	if cmdName == "" {
		cmdName = parts[0]
		args = parts[1:]
	}

	// resolve alias if exists
	if realCmdName, isAlias := r.aliases[cmdName]; isAlias {
		cmdName = realCmdName
	}

	// find the command. return CommandError if not found
	cmd, exists := r.commands[cmdName]
	if !exists {
		// If default command is set, execute it instead of returning error
		if r.defaultCmd != nil {
			return r.defaultCmd.Execute(args)
		}
		return &CommandError{
			Command: cmdName,
			Err:     fmt.Errorf("command not found"),
		}
	}

	// execute the command
	return cmd.Execute(args)
}

// create a Execute Command method that directly executes a command by name
func (r *Router) ExecuteCommand(cmdName string, args []string) error {
	// find the command. return CommandError if not found
	cmd, exists := r.commands[cmdName]
	if !exists {
		return &CommandError{
			Command: cmdName,
			Err:     fmt.Errorf("command not found"),
		}
	}

	// execute the command
	return cmd.Execute(args)
}

//--------------query methods------------------//

// create GetCommand
func (r *Router) GetCommand(cmdName string) (Command, bool) {
	// Resolve alias if exists
	if realCmdName, isAlias := r.aliases[cmdName]; isAlias {
		cmdName = realCmdName
	}
	// get command by name
	cmd, exists := r.commands[cmdName]
	return cmd, exists
}

// create ListCommands method.exclude aliases
func (r *Router) ListCommands() []Command {
	// return a slice of all commands, excluding aliases
	cmds := make([]Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}

// create HasCommand method that checks if a command exists
func (r *Router) HasCommand(cmdName string) bool {
	// Resolve alias if exists
	if realCmdName, isAlias := r.aliases[cmdName]; isAlias {
		cmdName = realCmdName
	}
	_, exists := r.commands[cmdName]
	return exists
}

// create GetAliases method that returns all aliases of a given command
func (r *Router) GetAliases() map[string]string {
	// return a copy of the aliases map
	aliases := make(map[string]string)
	for alias, cmdName := range r.aliases {
		aliases[alias] = cmdName
	}
	return aliases
}
