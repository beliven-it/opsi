// Package ui is the single place where opsi writes for humans: results go to
// stdout, messages about what is going on (warnings, errors, progress) go to
// stderr, and colors are used only when the destination is a terminal.
package ui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	gray   = "\033[90m"
)

// Where the output goes, replaceable by the tests
var (
	out    io.Writer = os.Stdout
	errOut io.Writer = os.Stderr

	// Whether each stream is a terminal that may get colors and progress
	outTerminal = isTerminal(os.Stdout)
	errTerminal = isTerminal(os.Stderr)
)

func isTerminal(f *os.File) bool {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}
	info, statErr := f.Stat()
	return statErr == nil && info.Mode()&os.ModeCharDevice != 0
}

func paint(enabled bool, color, text string) string {
	if !enabled {
		return text
	}
	return color + text + reset
}

// line writes a message ended by a newline, finishing a progress line first
func line(w io.Writer, terminal bool, color, symbol, format string, args ...any) {
	EndProgress()
	text := fmt.Sprintf(format, args...)
	if symbol != "" {
		text = paint(terminal, color, symbol) + " " + text
	}
	fmt.Fprintln(w, text)
}

// Print writes a plain line of the result
func Print(format string, args ...any) {
	line(out, outTerminal, "", "", format, args...)
}

// Success reports something done
func Success(format string, args ...any) {
	line(out, outTerminal, green, "✓", format, args...)
}

// Info reports a neutral fact, like a dry run or a nothing to do
func Info(format string, args ...any) {
	line(out, outTerminal, gray, "·", format, args...)
}

// Warn reports something that did not stop the command but needs attention
func Warn(format string, args ...any) {
	line(errOut, errTerminal, yellow, "!", format, args...)
}

// Error reports a failure
func Error(format string, args ...any) {
	line(errOut, errTerminal, red, "✗", format, args...)
}

// Fatal reports a failure and ends the command with a non zero exit code
func Fatal(failure error) {
	Error("%s", failure.Error())
	os.Exit(1)
}

// Section starts a group of lines, separated from the previous one
func Section(format string, args ...any) {
	EndProgress()
	fmt.Fprintf(out, "\n%s\n", paint(outTerminal, bold, fmt.Sprintf(format, args...)))
}

// Blank writes an empty line
func Blank() {
	EndProgress()
	fmt.Fprintln(out)
}

// Muted writes a secondary detail, like a path or an identifier
func Muted(format string, args ...any) {
	line(out, outTerminal, gray, "", "%s", paint(outTerminal, gray, fmt.Sprintf(format, args...)))
}

// Item writes an entry of a list, indented
func Item(format string, args ...any) {
	line(out, outTerminal, "", "", "  %s", fmt.Sprintf(format, args...))
}

// Table writes the rows aligned in columns, indented like the items. The last
// cell of a row is not padded, so it can carry a style.
func Table(rows [][]string) {
	EndProgress()
	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	for _, row := range rows {
		fmt.Fprintf(writer, "  %s\n", strings.Join(row, "\t"))
	}
	writer.Flush()
}

// Red, Yellow, Green, Dim and Bold return a text styled for stdout, to compose a line
func Red(text string) string    { return paint(outTerminal, red, text) }
func Yellow(text string) string { return paint(outTerminal, yellow, text) }
func Green(text string) string  { return paint(outTerminal, green, text) }
func Dim(text string) string    { return paint(outTerminal, gray, text) }
func Bold(text string) string   { return paint(outTerminal, bold, text) }

var progressActive bool

// Progress rewrites a status line on stderr. It is shown only on a terminal,
// so a pipe or a file never gets control characters.
func Progress(format string, args ...any) {
	if !errTerminal {
		return
	}
	fmt.Fprintf(errOut, "\r\033[K%s", paint(true, gray, fmt.Sprintf(format, args...)))
	progressActive = true
}

// EndProgress clears the status line, if any
func EndProgress() {
	if progressActive {
		fmt.Fprint(errOut, "\r\033[K")
		progressActive = false
	}
}

// Truncate shortens a text to max characters
func Truncate(text string, max int) string {
	runes := []rune(strings.TrimSpace(text))
	if len(runes) <= max {
		return string(runes)
	}
	return string(runes[:max-1]) + "…"
}
