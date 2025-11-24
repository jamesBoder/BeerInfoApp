package ui

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ToTitleCase converts a string to title case
func ToTitleCase(input string) string {
	caser := cases.Title(language.English)
	return caser.String(input)
}
