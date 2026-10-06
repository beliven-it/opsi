package ui

import (
	"bytes"
	"strings"
	"testing"
)

func capture(t *testing.T, terminal bool) (*bytes.Buffer, *bytes.Buffer) {
	t.Helper()
	var o, e bytes.Buffer
	oldOut, oldErr, oldOutT, oldErrT, oldActive := out, errOut, outTerminal, errTerminal, progressActive
	out, errOut, outTerminal, errTerminal, progressActive = &o, &e, terminal, terminal, false
	t.Cleanup(func() {
		out, errOut, outTerminal, errTerminal, progressActive = oldOut, oldErr, oldOutT, oldErrT, oldActive
	})
	return &o, &e
}

func TestStreams(t *testing.T) {
	o, e := capture(t, false)

	Success("created %s", "a")
	Info("nothing")
	Warn("careful")
	Error("broken")

	if got := o.String(); got != "✓ created a\n· nothing\n" {
		t.Errorf("stdout = %q", got)
	}
	if got := e.String(); got != "! careful\n✗ broken\n" {
		t.Errorf("stderr = %q", got)
	}
}

func TestNoColorOutsideTerminal(t *testing.T) {
	o, e := capture(t, false)

	Success("ok")
	Section("title")
	Error("ko")
	Progress("working")

	if strings.Contains(o.String()+e.String(), "\033") {
		t.Errorf("control characters outside a terminal: %q %q", o.String(), e.String())
	}
}

func TestColorsOnTerminal(t *testing.T) {
	o, e := capture(t, true)

	Success("ok")
	Error("ko")

	if !strings.Contains(o.String(), green+"✓"+reset) {
		t.Errorf("success not green: %q", o.String())
	}
	if !strings.Contains(e.String(), red+"✗"+reset) {
		t.Errorf("error not red: %q", e.String())
	}
}

func TestProgressIsClearedByNextLine(t *testing.T) {
	_, e := capture(t, true)

	Progress("step %d", 1)
	Progress("step %d", 2)
	Success("done")

	if progressActive {
		t.Error("progress still active after a line")
	}
	if !strings.Contains(e.String(), "step 2") || !strings.HasSuffix(e.String(), "\r\033[K") {
		t.Errorf("progress not cleared: %q", e.String())
	}
}

func TestTruncate(t *testing.T) {
	if got := Truncate("short", 10); got != "short" {
		t.Errorf("got %q", got)
	}
	if got := Truncate("abcdefghij", 5); got != "abcd…" {
		t.Errorf("got %q", got)
	}
}
