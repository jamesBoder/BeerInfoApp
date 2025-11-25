package cli

import (
	"strings"
	"testing"
)

// tests for parser.go file

// TestParseInput test, verify basic input parsing
func TestParseInput(t *testing.T) {
	input := "command arg1 arg2 arg3"
	cmdName, args := ParseInput(input)

	if cmdName != "command" {
		t.Errorf("Expected command name 'command', got '%s'", cmdName)
	}

	expectedArgs := []string{"arg1", "arg2", "arg3"}
	if len(args) != len(expectedArgs) {
		t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(args))
	}
	for i, arg := range expectedArgs {
		if args[i] != arg {
			t.Errorf("Expected arg %d to be '%s', got '%s'", i, arg, args[i])
		}
	}
}

// TestParseMultiWordInput test, verify multi-word command parsing
func TestParseMultiWordInput(t *testing.T) {
	multiWordCommands := []string{"multi word", "another command"}

	// test multi-word command
	input := "multi word arg1 arg2"
	cmdName, args := ParseMultiWordInput(input, multiWordCommands)

	if cmdName != "multi word" {
		t.Errorf("Expected command name 'multi word', got '%s'", cmdName)
	}

	expectedArgs := []string{"arg1", "arg2"}
	if len(args) != len(expectedArgs) {
		t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(args))
	}
	for i, arg := range expectedArgs {
		if args[i] != arg {
			t.Errorf("Expected arg %d to be '%s', got '%s'", i, arg, args[i])
		}
	}

	// test single-word fallback
	input = "singleword arg1 arg2"
	cmdName, args = ParseMultiWordInput(input, multiWordCommands)

	if cmdName != "singleword" {
		t.Errorf("Expected command name 'singleword', got '%s'", cmdName)
	}

	expectedArgs = []string{"arg1", "arg2"}
	if len(args) != len(expectedArgs) {
		t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(args))
	}
	for i, arg := range expectedArgs {
		if args[i] != arg {
			t.Errorf("Expected arg %d to be '%s', got '%s'", i, arg, args[i])
		}
	}
}

// TestParseMultiWordInput_Priority test, verify longer commands matched first
func TestParseMultiWordInput_Priority(t *testing.T) {
	multiWordCommands := []string{"multi word command", "multi word"}

	// input that matches both, should pick longer one
	input := "multi word command arg1 arg2"
	cmdName, args := ParseMultiWordInput(input, multiWordCommands)

	if cmdName != "multi word command" {
		t.Errorf("Expected command name 'multi word command', got '%s'", cmdName)
	}

	expectedArgs := []string{"arg1", "arg2"}
	if len(args) != len(expectedArgs) {
		t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(args))
	}
	for i, arg := range expectedArgs {
		if args[i] != arg {
			t.Errorf("Expected arg %d to be '%s', got '%s'", i, arg, args[i])
		}
	}
}

// TestSplitArgs test, verify argument splitting
func TestSplitArgs(t *testing.T) {
	input := "  arg1   arg2  arg3  "
	args := SplitArgs(input)

	expectedArgs := []string{"arg1", "arg2", "arg3"}
	if len(args) != len(expectedArgs) {
		t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(args))
	}
	for i, arg := range expectedArgs {
		if args[i] != arg {
			t.Errorf("Expected arg %d to be '%s', got '%s'", i, arg, args[i])
		}
	}
}

// TestQuoteAwareSplit test, verify quote handling
func TestQuoteAwareSplit(t *testing.T) {
	input := `arg1 "arg two with spaces" 'arg three' arg4`
	args := QuoteAwareSplit(input)

	expectedArgs := []string{"arg1", "arg two with spaces", "arg three", "arg4"}
	if len(args) != len(expectedArgs) {
		t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(args))
	}
	for i, arg := range expectedArgs {
		if args[i] != arg {
			t.Errorf("Expected arg %d to be '%s', got '%s'", i, arg, args[i])
		}
	}
}

// TestQuoteAwareSplit_EdgeCases test, verify edge cases for quotes
func TestQuoteAwareSplit_EdgeCases(t *testing.T) {
	// test unclosed quote
	input := `arg1 "arg two with spaces arg3`
	args := QuoteAwareSplit(input)

	expectedArgs := []string{"arg1", "arg two with spaces arg3"}
	if len(args) != len(expectedArgs) {
		t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(args))
	}
	for i, arg := range expectedArgs {
		if args[i] != arg {
			t.Errorf("Expected arg %d to be '%s', got '%s'", i, arg, args[i])
		}
	}

	// test empty input
	input = `   `
	args = QuoteAwareSplit(input)
	if len(args) != 0 {
		t.Errorf("Expected 0 args for empty input, got %d", len(args))
	}
}

// TestTrimINput test, verify input trimming
func TestTrimInput(t *testing.T) {
	input := "   some input string   "
	trimmed := strings.TrimSpace(input)
	expected := "some input string"
	if trimmed != expected {
		t.Errorf("Expected trimmed input to be '%s', got '%s'", expected, trimmed)
	}
}

// TestIsEmptyInput test, verify empty detection
func TestIsEmptyInput(t *testing.T) {
	input := "     "
	trimmed := strings.TrimSpace(input)
	if trimmed != "" {
		t.Errorf("Expected trimmed input to be empty, got '%s'", trimmed)
	}
}

// TestNormalizeInput test, verify input normalization
func TestNormalizeInput(t *testing.T) {
	input := "  COMMAND   Arg1   Arg2  "
	normalized := strings.ToLower(strings.TrimSpace(input))
	expected := "command   arg1   arg2"
	if normalized != expected {
		t.Errorf("Expected normalized input to be '%s', got '%s'", expected, normalized)
	}
}
