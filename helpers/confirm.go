package helpers

import (
	"bufio"
	"fmt"
	"opsi/helpers/ui"
	"os"
	"regexp"
	"strings"
)

func Confirm() {
	// Show the message
	fmt.Fprint(os.Stderr, "Are you sure to proceed? [y/N] ")

	// Start reader
	reader := bufio.NewReader(os.Stdin)

	// Read the STDIN and stop reading when
	// user press ENTER
	value, _ := reader.ReadString('\n')

	// Remove from the input string the \n character
	rgx := regexp.MustCompile(`\n`)
	value = rgx.ReplaceAllString(value, "")

	// Normalize the string in a lowercase format
	value = strings.ToLower(value)

	// Confirm if value is positive
	if value != "y" && value != "yes" {
		ui.Info("Ok, abort procedure")
		os.Exit(0)
	}
}
