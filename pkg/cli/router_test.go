package cli

import (
	"testing"
)

// tests for router.go file

// TestNewRouter returns non-nil pointer, cmd map is init, alias map is init, multiWordCmds slice is init
func TestNewRouter(t *testing.T) {
	router := NewRouter()
	if router == nil {
		t.Fatal("Expected non-nil Router pointer")
	}
	if router.commands == nil {
		t.Error("Expected commands map to be initialized")
	}
	if router.aliases == nil {
		t.Error("Expected aliases map to be initialized")
	}
	if router.multiWordCmds == nil {
		t.Error("Expected multiWordCmds slice to be initialized")
	}
}

// TestRouter_Register test registers cmd, cmd is retrievable, returns error for dups, returns error for empty name, handles multiple cmds
func TestRouter_RegisterCommand(t *testing.T) {
	router := NewRouter()

	// create test command
	handler := func(args []string) error {
		return nil
	}
	cmd := NewBaseCommand("test", "A test command", "test [args]", handler)

	// register command
	err := router.RegisterCommand(cmd)
	if err != nil {
		t.Fatalf("Expected no error registering command, got: %v", err)
	}

	// verify command is retrievable
	registeredCmd, exists := router.commands["test"]
	if !exists {
		t.Fatal("Expected command 'test' to be registered")
	}
	if registeredCmd.Name() != "test" {
		t.Errorf("Expected registered command name 'test', got '%s'", registeredCmd.Name())
	}

	// attempt to register duplicate command
	err = router.RegisterCommand(cmd)
	if err == nil {
		t.Fatal("Expected error when registering duplicate command, got nil")
	}

	// attempt to register command with empty name
	emptyNameCmd := NewBaseCommand("", "No name command", "no-name", handler)
	err = router.RegisterCommand(emptyNameCmd)
	if err == nil {
		t.Fatal("Expected error when registering command with empty name, got nil")
	}

	// register another command
	anotherCmd := NewBaseCommand("another", "Another command", "another [args]", handler)
	err = router.RegisterCommand(anotherCmd)
	if err != nil {
		t.Fatalf("Expected no error registering another command, got: %v", err)
	}

	// verify both commands exist
	if len(router.commands) != 2 {
		t.Errorf("Expected 2 registered commands, got %d", len(router.commands))
	}
}

// TestRouter_RegisterCommand_Duplicate tests if second registration fails, first cmd remains, error msg is clear
func TestRouter_RegisterCommand_Duplicate(t *testing.T) {
	router := NewRouter()

	// create test command
	handler := func(args []string) error {
		return nil
	}
	cmd := NewBaseCommand("test", "A test command", "test [args]", handler)

	// register command first time
	err := router.RegisterCommand(cmd)
	if err != nil {
		t.Fatalf("Expected no error registering command, got: %v", err)
	}

	// register command second time (duplicate)
	err = router.RegisterCommand(cmd)
	if err == nil {
		t.Fatal("Expected error when registering duplicate command, got nil")
	}

	// verify error message is clear
	expectedErrMsg := "command already registered: test"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}

	// verify only one instance of the command exists
	if len(router.commands) != 1 {
		t.Errorf("Expected 1 registered command, got %d", len(router.commands))
	}
}

// TestRouter_RegisterCommand_EmptyName tests if registration fails, error msg is clear, no cmd is added
func TestRouter_RegisterCommand_EmptyName(t *testing.T) {
	router := NewRouter()

	// create test command with empty name
	handler := func(args []string) error {
		return nil
	}
	cmd := NewBaseCommand("", "No name command", "no-name", handler)

	// attempt to register command
	err := router.RegisterCommand(cmd)
	if err == nil {
		t.Fatal("Expected error when registering command with empty name, got nil")
	}

	// verify error message is clear
	expectedErrMsg := "command name cannot be empty"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}

	// verify no command was added
	if len(router.commands) != 0 {
		t.Errorf("Expected 0 registered commands, got %d", len(router.commands))
	}
}

// TestRouter_RegistAlias , creates alias, points to correct cmd, returns err if cmd doesnt exist, returns err for dup alias, handles multiple aliases
func TestRouter_RegisterAlias(t *testing.T) {
	router := NewRouter()

	// create test command
	handler := func(args []string) error {
		return nil
	}

	// register command
	cmd := NewBaseCommand("test", "A test command", "test [args]", handler)
	err := router.RegisterCommand(cmd)
	if err != nil {
		t.Fatalf("Expected no error registering command, got: %v", err)
	}

	// register alias
	err = router.RegisterAlias("t", "test")
	if err != nil {
		t.Fatalf("Expected no error registering alias, got: %v", err)
	}

	// verify alias points to correct command
	cmdName, exists := router.aliases["t"]
	if !exists {
		t.Fatal("Expected alias 't' to be registered")
	}
	if cmdName != "test" {
		t.Errorf("Expected alias 't' to point to command 'test', got '%s'", cmdName)
	}

	// attempt to register alias for non-existent command
	err = router.RegisterAlias("x", "nonexistent")
	if err == nil {
		t.Fatal("Expected error when registering alias for non-existent command, got nil")
	}

	// attempt to register duplicate alias
	err = router.RegisterAlias("t", "test")
	if err != nil {
		t.Fatalf("Expected no error when registering duplicate alias, got: %v", err)
	}

	// register another alias
	err = router.RegisterAlias("testcmd", "test")
	if err != nil {
		t.Fatalf("Expected no error registering another alias, got: %v", err)
	}

	// verify both aliases exist
	if len(router.aliases) != 2 {
		t.Errorf("Expected 2 registered aliases, got %d", len(router.aliases))
	}
}

// TestRouter_RegisterMultiWordCommand , marks cmd as multi-word, returns err if cmd doesnt exist, handles multiple multi-word cmds, order matters(longer first)
func TestRouter_RegisterMultiWordCommand(t *testing.T) {
	router := NewRouter()

	// create test commands
	handler := func(args []string) error {
		return nil
	}

	// register commands
	cmd1 := NewBaseCommand("multi word", "A multi word command", "multi word [args]", handler)
	err := router.RegisterCommand(cmd1)
	if err != nil {
		t.Fatalf("Expected no error registering command, got: %v", err)
	}

	// register multi-word command
	mwCmd := NewBaseCommand("another multi word", "Another multi word command", "another multi word [args]", handler)

	err = router.RegisterCommand(mwCmd)
	if err != nil {
		t.Fatalf("Expected no error registering command, got: %v", err)
	}

	// mark first command as multi-word
	err = router.RegisterMultiWordCommand("multi word")
	if err != nil {
		t.Fatalf("Expected no error registering multi-word command, got: %v", err)
	}

	// mark second command as multi-word
	err = router.RegisterMultiWordCommand("another multi word")
	if err != nil {
		t.Fatalf("Expected no error registering multi-word command, got: %v", err)
	}

	// verify both multi-word commands are registered
	if len(router.multiWordCmds) != 2 {
		t.Errorf("Expected 2 registered multi-word commands, got %d", len(router.multiWordCmds))
	}

	// attempt to register multi-word command that doesn't exist
	err = router.RegisterMultiWordCommand("nonexistent command")
	if err == nil {
		t.Fatal("Expected error when registering non-existent multi-word command, got nil")
	}
}

// TestRouter_Execute_SimpleCommand tests cmd with no args, cmd with single arg, cmd with multiple args, cmd execution is called correctly
func TestRouter_Execute_SimpleCommand(t *testing.T) {
	router := NewRouter()

	// create test command
	executed := false
	var receivedArgs []string
	handler := func(args []string) error {
		executed = true
		receivedArgs = args
		return nil
	}
	cmd := NewBaseCommand("test", "A test command", "test [args]", handler)

	// register command
	err := router.RegisterCommand(cmd)
	if err != nil {
		t.Fatalf("Expected no error registering command, got: %v", err)
	}

	// execute command with no args
	executed = false
	receivedArgs = nil
	err = router.Execute("test")
	if err != nil {
		t.Fatalf("Expected no error executing command, got: %v", err)
	}
	if !executed {
		t.Fatal("Expected command to be executed")
	}
	if len(receivedArgs) != 0 {
		t.Errorf("Expected 0 args, got %d", len(receivedArgs))
	}

	// execute command with single arg
	executed = false
	receivedArgs = nil
	err = router.Execute("test arg1")
	if err != nil {
		t.Fatalf("Expected no error executing command, got: %v", err)
	}
	if !executed {
		t.Fatal("Expected command to be executed")
	}
	if len(receivedArgs) != 1 || receivedArgs[0] != "arg1" {
		t.Errorf("Expected args ['arg1'], got %v", receivedArgs)
	}

	// execute command with multiple args
	executed = false
	receivedArgs = nil
	err = router.Execute("test arg1 arg2 arg3")
	if err != nil {
		t.Fatalf("Expected no error executing command, got: %v", err)
	}
	if !executed {
		t.Fatal("Expected command to be executed")
	}
	if len(receivedArgs) != 3 || receivedArgs[0] != "arg1" || receivedArgs[1] != "arg2" || receivedArgs[2] != "arg3" {
		t.Errorf("Expected args ['arg1', 'arg2', 'arg3'], got %v", receivedArgs)
	}
}

// TestRouter_Execute_MultiWordCommand tests multi-word cmd with no args, with single arg, with multiple args, cmd execution is called correctly. falls back to single-word cmd
func TestRouter_Execute_MultiWordCommand(t *testing.T) {
	router := NewRouter()

	// create test multi-word command
	executed := false
	var receivedArgs []string
	handler := func(args []string) error {
		executed = true
		receivedArgs = args
		return nil
	}
	mwCmd := NewBaseCommand("multi word", "A multi word command", "multi word [args]", handler)

	// register command
	err := router.RegisterCommand(mwCmd)
	if err != nil {
		t.Fatalf("Expected no error registering command, got: %v", err)
	}

	// mark as multi-word command
	err = router.RegisterMultiWordCommand("multi word")
	if err != nil {
		t.Fatalf("Expected no error registering multi-word command, got: %v", err)
	}

	// execute multi-word command with no args
	executed = false
	receivedArgs = nil
	err = router.Execute("multi word")
	if err != nil {
		t.Fatalf("Expected no error executing command, got: %v", err)
	}
	if !executed {
		t.Fatal("Expected command to be executed")
	}
	if len(receivedArgs) != 0 {
		t.Errorf("Expected 0 args, got %d", len(receivedArgs))
	}

	// execute multi-word command with single arg
	executed = false
	receivedArgs = nil
	err = router.Execute("multi word arg1")
	if err != nil {
		t.Fatalf("Expected no error executing command, got: %v", err)
	}
	if !executed {
		t.Fatal("Expected command to be executed")
	}
	if len(receivedArgs) != 1 || receivedArgs[0] != "arg1" {
		t.Errorf("Expected args ['arg1'], got %v", receivedArgs)
	}

	// execute multi-word command with multiple args
	executed = false
	receivedArgs = nil
	err = router.Execute("multi word arg1 arg2 arg3")
	if err != nil {
		t.Fatalf("Expected no error executing command, got: %v", err)
	}
	if !executed {
		t.Fatal("Expected command to be executed")
	}
	if len(receivedArgs) != 3 || receivedArgs[0] != "arg1" || receivedArgs[1] != "arg2" || receivedArgs[2] != "arg3" {
		t.Errorf("Expected args ['arg1', 'arg2', 'arg3'], got %v", receivedArgs)
	}

	// execute non-matching multi-word command (falls back to single-word)
	executed = false
	receivedArgs = nil
	err = router.Execute("multi arg1 arg2")
	if err == nil {
		t.Fatal("Expected error executing non-matching multi-word command, got nil")
	}
}

// TestRouter_Execute_WithAlias test,  alias exec target cmd, args are passed through, multiple aliases work, alias doesn't appear in teh cmd list
func TestRouter_Execute_WithAlias(t *testing.T) {
	router := NewRouter()

	// create test command
	executed := false
	var receivedArgs []string
	handler := func(args []string) error {
		executed = true
		receivedArgs = args
		return nil
	}
	cmd := NewBaseCommand("test", "A test command", "test [args]", handler)

	// register command
	err := router.RegisterCommand(cmd)
	if err != nil {
		t.Fatalf("Expected no error registering command, got: %v", err)
	}

	// register alias
	err = router.RegisterAlias("t", "test")
	if err != nil {
		t.Fatalf("Expected no error registering alias, got: %v", err)
	}

	// execute command via alias with args
	executed = false
	receivedArgs = nil
	err = router.Execute("t arg1 arg2")
	if err != nil {
		t.Fatalf("Expected no error executing command via alias, got: %v", err)
	}

	if !executed {
		t.Fatal("Expected command to be executed via alias")
	}
	if len(receivedArgs) != 2 || receivedArgs[0] != "arg1" || receivedArgs[1] != "arg2" {
		t.Errorf("Expected args ['arg1', 'arg2'], got %v", receivedArgs)
	}

	// verify alias does not appear in command list
	if _, exists := router.commands["t"]; exists {
		t.Error("Expected alias 't' to not appear in command list")
	}
}

// TestRouter_Execute_NotFound tests, unknown cmd returns err, err is CommandError type, error msg includes cmd name, defaul cmd is called if set
func TestRouter_Execute_NotFound(t *testing.T) {
	router := NewRouter()

	// execute unknown command
	err := router.Execute("unknowncmd arg1 arg2")
	if err == nil {
		t.Fatal("Expected error executing unknown command, got nil")
	}

	// verify error is of type CommandError
	cmdErr, ok := err.(*CommandError)
	if !ok {
		t.Fatalf("Expected error of type CommandError, got %T", err)
	}

	// verify error message includes command name
	expectedCmdName := "unknowncmd"
	if cmdErr.Command != expectedCmdName {
		t.Errorf("Expected CommandError.Command to be '%s', got '%s'", expectedCmdName, cmdErr.Command)
	}

	// set default command
	executed := false
	handler := func(args []string) error {
		executed = true
		return nil
	}
	defaultCmd := NewBaseCommand("default", "Default command", "default [args]", handler)
	router.SetDefaultCommand(defaultCmd)

	// execute unknown command again
	executed = false
	err = router.Execute("unknowncmd arg1 arg2")
	if err != nil {
		t.Fatalf("Expected no error executing unknown command with default set, got: %v", err)
	}

	// verify default command was executed
	if !executed {
		t.Fatal("Expected default command to be executed for unknown command")
	}
}

// TestRouter_Execute_EmptyInput tests, empty string returns err, whitespace only returns err, default cmd is called if set
func TestRouter_Execute_EmptyInput(t *testing.T) {
	router := NewRouter()

	// execute empty input
	err := router.Execute("")
	if err == nil {
		t.Fatal("Expected error executing empty input, got nil")
	}

	// execute whitespace only input
	err = router.Execute("    ")
	if err == nil {
		t.Fatal("Expected error executing whitespace only input, got nil")
	}

	// set default command
	executed := false
	handler := func(args []string) error {
		executed = true
		return nil
	}
	defaultCmd := NewBaseCommand("default", "Default command", "default [args]", handler)
	router.SetDefaultCommand(defaultCmd)

	// execute empty input again
	executed = false
	err = router.Execute("   ")
	if err != nil {
		t.Fatalf("Expected no error executing whitespace input with default set, got: %v", err)
	}

	// verify default command was executed
	if !executed {
		t.Fatal("Expected default command to be executed for empty input")
	}
}

// TestRouter_ExecuteCommand test , executes cmd by name, passes args correctly, returns err for unknown cmd, doesnt resolve aliases

func TestRouter_ExecuteCommand(t *testing.T) {
	router := NewRouter()

	// create test command
	executed := false
	var receivedArgs []string
	handler := func(args []string) error {
		executed = true
		receivedArgs = args
		return nil
	}
	cmd := NewBaseCommand("test", "A test command", "test [args]", handler)

	// register command
	err := router.RegisterCommand(cmd)
	if err != nil {
		t.Fatalf("Expected no error registering command, got: %v", err)
	}

	// execute command by name with args
	executed = false
	receivedArgs = nil
	err = router.ExecuteCommand("test", []string{"arg1", "arg2"})
	if err != nil {
		t.Fatalf("Expected no error executing command by name, got: %v", err)
	}
	if !executed {
		t.Fatal("Expected command to be executed by name")
	}
	if len(receivedArgs) != 2 || receivedArgs[0] != "arg1" || receivedArgs[1] != "arg2" {
		t.Errorf("Expected args ['arg1', 'arg2'], got %v", receivedArgs)
	}

	// attempt to execute unknown command by name
	err = router.ExecuteCommand("unknowncmd", []string{})
	if err == nil {
		t.Fatal("Expected error executing unknown command by name, got nil")
	}
}

// TestRouter_GetCommand test, returns cmd if exists, returns false if not exists, resolves aliases, case sensitive
func TestRouter_GetCommand(t *testing.T) {
	router := NewRouter()

	// create test command
	handler := func(args []string) error {
		return nil
	}
	cmd := NewBaseCommand("test", "A test command", "test [args]", handler)

	// register command
	err := router.RegisterCommand(cmd)
	if err != nil {
		t.Fatalf("Expected no error registering command, got: %v", err)
	}

	// register alias
	err = router.RegisterAlias("t", "test")
	if err != nil {
		t.Fatalf("Expected no error registering alias, got: %v", err)
	}

	// get command by name
	retrievedCmd, exists := router.GetCommand("test")
	if !exists {
		t.Fatal("Expected command 'test' to exist")
	}
	if retrievedCmd.Name() != "test" {
		t.Errorf("Expected retrieved command name 'test', got '%s'", retrievedCmd.Name())
	}

	// get command by alias
	retrievedCmd, exists = router.GetCommand("t")
	if !exists {
		t.Fatal("Expected command for alias 't' to exist")
	}
	if retrievedCmd.Name() != "test" {
		t.Errorf("Expected retrieved command name 'test' for alias 't', got '%s'", retrievedCmd.Name())
	}

	// attempt to get non-existent command
	_, exists = router.GetCommand("unknowncmd")
	if exists {
		t.Fatal("Expected command 'unknowncmd' to not exist")
	}

	// attempt to get command with different case
	_, exists = router.GetCommand("Test")
	if exists {
		t.Fatal("Expected command 'Test' (case sensitive) to not exist")
	}
}

// TestRouter_ListCommands test, returns all cmds, doesn't include aliases, order doesn't matter, empty router returns empty slice
func TestRouter_ListCommands(t *testing.T) {
	router := NewRouter()

	// create test commands
	handler := func(args []string) error {
		return nil
	}
	cmd1 := NewBaseCommand("cmd1", "First command", "cmd1 [args]", handler)
	cmd2 := NewBaseCommand("cmd2", "Second command", "cmd2 [args]", handler)

	// register commands
	err := router.RegisterCommand(cmd1)
	if err != nil {
		t.Fatalf("Expected no error registering command, got: %v", err)
	}
	err = router.RegisterCommand(cmd2)
	if err != nil {
		t.Fatalf("Expected no error registering command, got: %v", err)
	}

	// register alias
	err = router.RegisterAlias("c1", "cmd1")
	if err != nil {
		t.Fatalf("Expected no error registering alias, got: %v", err)
	}

	// list commands
	cmds := router.ListCommands()
	if len(cmds) != 2 {
		t.Fatalf("Expected 2 commands in list, got %d", len(cmds))
	}

	// verify commands in list
	foundCmd1 := false
	foundCmd2 := false
	for _, c := range cmds {
		if c.Name() == "cmd1" {
			foundCmd1 = true
		} else if c.Name() == "cmd2" {
			foundCmd2 = true
		}
	}
	if !foundCmd1 {
		t.Error("Expected command 'cmd1' to be in the list")
	}
	if !foundCmd2 {
		t.Error("Expected command 'cmd2' to be in the list")
	}

	// test empty router
	emptyRouter := NewRouter()
	cmds = emptyRouter.ListCommands()
	if len(cmds) != 0 {
		t.Fatalf("Expected 0 commands in empty router list, got %d", len(cmds))
	}
}

// TestRouter_HasCommand tests, returns true for existing command, returns false for non-existing cmd, resolves aliases, case sensitive
func TestRouter_HasCommand(t *testing.T) {
	router := NewRouter()

	// create test command
	handler := func(args []string) error {
		return nil
	}
	cmd := NewBaseCommand("test", "A test command", "test [args]", handler)

	// register command
	err := router.RegisterCommand(cmd)
	if err != nil {
		t.Fatalf("Expected no error registering command, got: %v", err)
	}

	// register alias
	err = router.RegisterAlias("t", "test")
	if err != nil {
		t.Fatalf("Expected no error registering alias, got: %v", err)
	}

	// check existing command
	if !router.HasCommand("test") {
		t.Fatal("Expected HasCommand to return true for existing command 'test'")
	}

	// check existing command via alias
	if !router.HasCommand("t") {
		t.Fatal("Expected HasCommand to return true for alias 't'")
	}

	// check non-existing command
	if router.HasCommand("unknowncmd") {
		t.Fatal("Expected HasCommand to return false for non-existing command 'unknowncmd'")
	}

	// check command with different case
	if router.HasCommand("Test") {
		t.Fatal("Expected HasCommand to return false for command 'Test' (case sensitive)")
	}
}

// TestRouter_GetAlias test, returns all aliases, map is copy, empty if no aliases, correct cmd mapping

func TestRouter_GetAliases(t *testing.T) {
	router := NewRouter()
	// create test command
	handler := func(args []string) error {
		return nil
	}
	cmd := NewBaseCommand("test", "A test command", "test [args]", handler)

	// register command
	err := router.RegisterCommand(cmd)
	if err != nil {
		t.Fatalf("Expected no error registering command, got: %v", err)
	}

	// register aliases
	err = router.RegisterAlias("t", "test")
	if err != nil {
		t.Fatalf("Expected no error registering alias, got: %v", err)
	}
	err = router.RegisterAlias("testcmd", "test")
	if err != nil {
		t.Fatalf("Expected no error registering alias, got: %v", err)
	}

	// get aliases
	aliases := router.GetAliases()
	if len(aliases) != 2 {
		t.Fatalf("Expected 2 aliases, got %d", len(aliases))
	}

	// verify aliases map
	if cmdName, exists := aliases["t"]; !exists || cmdName != "test" {
		t.Error("Expected alias 't' to map to command 'test'")
	}
	if cmdName, exists := aliases["testcmd"]; !exists || cmdName != "test" {
		t.Error("Expected alias 'testcmd' to map to command 'test'")
	}

	// modify returned map and verify original is unchanged
	aliases["newalias"] = "newcmd"
	if _, exists := router.aliases["newalias"]; exists {
		t.Error("Expected original aliases map to be unchanged when modifying returned map")
	}
}

// TestRouter_SetDefaultCommand test, default command is set, called on empty input, called on unknown cmd, can be nil
func TestRouter_SetDefaultCommand(t *testing.T) {
	router := NewRouter()

	// create test default command
	executed := false
	handler := func(args []string) error {
		executed = true
		return nil
	}
	defaultCmd := NewBaseCommand("default", "Default command", "default [args]", handler)

	// set default command
	router.SetDefaultCommand(defaultCmd)

	// execute empty input
	executed = false
	err := router.Execute("")
	if err != nil {
		t.Fatalf("Expected no error executing empty input with default set, got: %v", err)
	}
	if !executed {
		t.Fatal("Expected default command to be executed for empty input")
	}

	// execute unknown command
	executed = false
	err = router.Execute("unknowncmd arg1 arg2")
	if err != nil {
		t.Fatalf("Expected no error executing unknown command with default set, got: %v", err)
	}
	if !executed {
		t.Fatal("Expected default command to be executed for unknown command")
	}

	// set default command to nil
	router.SetDefaultCommand(nil)

	// execute empty input again
	err = router.Execute("")
	if err == nil {
		t.Fatal("Expected error executing empty input with nil default, got nil")
	}

	// execute unknown command again
	err = router.Execute("anotherunknown arg1 arg2")
	if err == nil {
		t.Fatal("Expected error executing unknown command with nil default, got nil")
	}
}
