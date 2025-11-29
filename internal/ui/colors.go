package ui

import (
	"fmt"

	"github.com/fatih/color"
)

// create Error() color msg
func Error(format string, args ...interface{}) string {
	msg := fmt.Sprintf(format, args...)
	return color.RedString("❌ " + msg)
}

// create Success() color msg
func Success(format string, args ...interface{}) string {
	msg := fmt.Sprintf(format, args...)
	return color.GreenString("✅ " + msg)
}

// create Warning() color msg
func Warning(format string, args ...interface{}) string {
	msg := fmt.Sprintf(format, args...)
	return color.YellowString("⚠️ " + msg)
}

// create Info() color msg
func Info(format string, args ...interface{}) string {
	msg := fmt.Sprintf(format, args...)
	return color.CyanString("ℹ️ " + msg)
}

// create Tip() color msg
func Tip(format string, args ...interface{}) string {
	msg := fmt.Sprintf(format, args...)
	return color.MagentaString("💡 " + msg)
}

// create Example() color msg
func Example(format string, args ...interface{}) string {
	msg := fmt.Sprintf(format, args...)
	return color.BlueString("📖 " + msg)
}

// create a Prompt() color msg
func Prompt(format string, args ...interface{}) string {
	msg := fmt.Sprintf(format, args...)
	return color.New(color.FgHiGreen, color.Bold).Sprintf("➡️  %s", msg)
}

// create a FieldLabel() color msg
func FieldLabel(format string, args ...interface{}) string {
	msg := fmt.Sprintf(format, args...)
	return color.New(color.FgHiMagenta, color.Bold).Sprintf("%s:", msg)
}

// create Header() color msg
func Header(format string, args ...interface{}) string {
	msg := fmt.Sprintf(format, args...)
	return color.New(color.FgHiWhite, color.Bold).Sprintf("🔷 %s", msg)
}

// create Divider() color msg
func Divider() string {
	return color.New(color.FgCyan).Sprintf("-------------------------------------------")
}

// create a Separator() color msg for the cyan dashes between fields
func Separator() string {
	return color.New(color.FgCyan).Sprintf("──────────────────")
}

// create a ****** line color msg cyan
func StarDivider() string {
	return color.New(color.FgCyan).Sprintf("*****************************************************************************************************************")
}

// create BeerNumber() color msg
func BeerNumber(num int) string {
	return color.New(color.FgHiYellow, color.Bold).Sprintf("🍺 Beer #%d", num)
}

// CyanBold helper function to return cyan bold string
func CyanBold(text string) string {
	return color.New(color.FgCyan, color.Bold).Sprintf(text)
}

// create SearchNumber() color msg
func SearchNumber(num int) string {
	return color.New(color.FgHiYellow, color.Bold).Sprintf("🔍 Search #%d", num)
}
