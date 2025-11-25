package cli

import (
	"strings"
)

// Parse user input into command names and arguments, handling various input formats.

// parsing function for single word commands
func ParseInput(input string) (commandName string, args []string) {
	// trim whitespace
	input = strings.TrimSpace(input)

	// split input into parts
	parts := strings.Fields(input)

	// if no input, return empty command
	if len(parts) == 0 {
		return "", []string{}
	}

	// first part is command name
	commandName = parts[0]

	// remaining parts are arguments
	if len(parts) > 1 {
		args = parts[1:]
	} else {
		args = []string{}
	}

	return commandName, args
}

// create a parsing function for multi-word commands
func ParseMultiWordInput(input string, multiWordCommands []string) (commandName string, args []string) {
	// trim whitespace
	input = strings.TrimSpace(input)

	// check if input starts with any multi-word command
	for _, mwCmd := range multiWordCommands {
		if strings.HasPrefix(input, mwCmd) {
			// found a matching multi-word command
			commandName = mwCmd
			// remaining text after becomes arguments
			remaining := strings.TrimSpace(strings.TrimPrefix(input, mwCmd))
			if remaining != "" {
				args = strings.Fields(remaining)
			} else {
				args = []string{}
			}
			return commandName, args
		}
	}

	// if no multi-word command matched, fall back to single word parsing
	return ParseInput(input)

}

// create a splitArgs function to split args by whitespace, filter empty strings, split on spaces
func SplitArgs(input string) []string {

	// trim whitespace
	input = strings.TrimSpace(input)

	// split on spaces
	parts := strings.Split(input, " ")

	// filter out empty strings
	var args []string
	for _, part := range parts {
		if part != "" {
			args = append(args, part)
		}
	}
	return args
}

// create a function that splits arguments respecting quoted strings
func QuoteAwareSplit(input string) []string {
	var args []string
	var currentArg strings.Builder
	inQuotes := false
	var quoteChar rune

	for _, char := range input {
		switch char {
		case ' ':
			if inQuotes {
				currentArg.WriteRune(char)
			} else if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
		case '"', '\'':
			if !inQuotes {
				// Start of quoted section
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				// End of quoted section (matching quote)
				inQuotes = false
				quoteChar = 0
			} else {
				// Different quote character while in quotes
				currentArg.WriteRune(char)
			}
		default:
			currentArg.WriteRune(char)
		}
	}

	// add last arg if exists
	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	return args
}

// -------------utitily functions------------------//

// create a trimInput function to trim leading and trailing whitespace
func TrimInput(input string) string {
	input = strings.TrimSpace(input)
	return input
}

// create a isEmptyInput function to check if input is empty or only whitespace
func IsEmptyInput(input string) bool {
	// check if input is empty
	if input == "" || len(input) == 0 {
		return true
	}
	// check if input is only whitespace
	if strings.TrimSpace(input) == "" {
		return true
	}
	return false
}

// create a normalizeInput function to convert input to lowercase
func NormalizeInput(input string) string {
	// trim whitespace
	input = strings.TrimSpace(input)
	// convert to lowercase
	input = strings.ToLower(input)
	return input
}
