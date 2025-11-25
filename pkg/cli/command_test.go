package cli

import (
	"fmt"
	"testing"
)

// tests for command.go file

// Helper function to create a test command (not a test itself)
func newTestCommand(t *testing.T) *BaseCommand {
	t.Helper() // marks this as a helper function
	handler := func(args []string) error {
		return nil
	}
	return NewBaseCommand("test", "A test command", "test [args]", handler)
}

// TestNewBaseCommand tests the NewBaseCommand function. creates command with all fields. returns non-nil pointer
// fields are set correctly
func TestNewBaseCommand(t *testing.T) {
	handler := func(args []string) error {
		return nil
	}
	cmd := NewBaseCommand("test", "A test command", "test [args]", handler)
	if cmd == nil {
		t.Fatal("Expected non-nil BaseCommand pointer")
	}
	if cmd.name != "test" {
		t.Errorf("Expected name 'test', got '%s'", cmd.name)
	}
	if cmd.description != "A test command" {
		t.Errorf("Expected description 'A test command', got '%s'", cmd.description)
	}
	if cmd.usage != "test [args]" {
		t.Errorf("Expected usage 'test [args]', got '%s'", cmd.usage)
	}
	if cmd.handler == nil {
		t.Error("Expected non-nil handler function")
	}
}

// TestBaseCommand_Name tests the Name method of BaseCommand
func TestBaseCommand_Name(t *testing.T) {
	cmd := newTestCommand(t)
	if cmd.Name() != "test" {
		t.Errorf("Expected command name 'test', got '%s'", cmd.Name())
	}
}

// TestBaseCommand_Description tests the Description method of BaseCommand
func TestBaseCommand_Description(t *testing.T) {
	cmd := newTestCommand(t)
	if cmd.Description() != "A test command" {
		t.Errorf("Expected command description 'A test command', got '%s'", cmd.Description())
	}
}

// TestBaseCommand_Usage tests the Usage method of BaseCommand
func TestBaseCommand_Usage(t *testing.T) {
	cmd := newTestCommand(t)
	if cmd.Usage() != "test [args]" {
		t.Errorf("Expected command usage 'test [args]', got '%s'", cmd.Usage())
	}
}

// TestBaseCommand_Execute tests the Execute method of BaseCommand
func TestBaseCommand_Execute(t *testing.T) {
	executed := false
	cmd := NewBaseCommand("test", "A test command", "test [args]", func(args []string) error {
		executed = true
		return nil
	})
	err := cmd.Execute([]string{})
	if err != nil {
		t.Errorf("Expected no error from Execute, got '%v'", err)
	}
	if !executed {
		t.Errorf("Expected command handler to be executed")
	}
}

// TesstBaseCommand_ExecuteWithArgs. verify Execute() handles various args
func TestBaseCommand_ExecuteWithArgs(t *testing.T) {
	var receivedArgs []string
	cmd := NewBaseCommand("test", "A test command", "test [args]", func(args []string) error {
		receivedArgs = args
		return nil
	})
	testArgs := []string{"arg1", "arg2", "arg3"}
	err := cmd.Execute(testArgs)
	if err != nil {
		t.Errorf("Expected no error from Execute, got '%v'", err)
	}
	if len(receivedArgs) != len(testArgs) {
		t.Errorf("Expected %d args, got %d", len(testArgs), len(receivedArgs))
	}
	for i, arg := range testArgs {
		if receivedArgs[i] != arg {
			t.Errorf("Expected arg %d to be '%s', got '%s'", i, arg, receivedArgs[i])
		}
	}
}

// TestCommandError_Error tests the Error method of CommandError
func TestCommandError_Error(t *testing.T) {
	cmdErr := &CommandError{
		Command: "test",
		Err:     fmt.Errorf("an error occurred"),
	}
	expectedMsg := "command 'test': an error occurred"
	if cmdErr.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, cmdErr.Error())
	}
}

// TestCommandError_Unwrap tests the Unwrap method of CommandError
func TestCommandError_Unwrap(t *testing.T) {
	innerErr := fmt.Errorf("inner error")
	cmdErr := &CommandError{
		Command: "test",
		Err:     innerErr,
	}
	if cmdErr.Unwrap() != innerErr {
		t.Errorf("Expected unwrapped error to be '%v', got '%v'", innerErr, cmdErr.Unwrap())
	}
}

// End of tests for command.go file
