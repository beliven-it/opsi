package main

import (
	"opsi/helpers"
	"os"
	"testing"
)

func TestConfigInitWritesTheTemplateWhenMissing(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := helpers.ConfigInit(configTemplate, "/.config/opsi/config.yml"); err != nil {
		t.Fatal(err)
	}

	written, err := os.ReadFile(home + "/.config/opsi/config.yml")
	if err != nil {
		t.Fatalf("the config should have been created: %v", err)
	}

	expected, _ := configTemplate.ReadFile("config/template.yml")
	if string(written) != string(expected) {
		t.Errorf("the config should match the shipped template, got %q", written)
	}

	// The file holds api tokens, so no other user on the machine may read it
	info, err := os.Stat(home + "/.config/opsi/config.yml")
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != 0600 {
		t.Errorf("the config should be 0600, got %o", info.Mode().Perm())
	}

	dir, err := os.Stat(home + "/.config/opsi")
	if err != nil {
		t.Fatal(err)
	}

	if dir.Mode().Perm() != 0700 {
		t.Errorf("the config directory should be 0700, got %o", dir.Mode().Perm())
	}
}

func TestConfigInitLeavesAnExistingConfigAlone(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	path := home + "/.config/opsi/config.yml"
	if err := os.MkdirAll(home+"/.config/opsi", 0755); err != nil {
		t.Fatal(err)
	}

	mine := "gitlab:\n  token: mine\n"
	if err := os.WriteFile(path, []byte(mine), 0644); err != nil {
		t.Fatal(err)
	}

	if err := helpers.ConfigInit(configTemplate, "/.config/opsi/config.yml"); err != nil {
		t.Fatal(err)
	}

	// This is why a key added to the template never reaches an existing user
	kept, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(kept) != mine {
		t.Errorf("an existing config must not be overwritten, got %q", kept)
	}
}
