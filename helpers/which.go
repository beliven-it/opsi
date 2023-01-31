package helpers

import (
	"bytes"
	"errors"
	"os/exec"
)

func Which(command string) string {
	response, err := Exec("which", command)
	if err != nil {
		return ""
	}

	return string(response)
}

func Exec(command string, args ...string) ([]byte, error) {
	var stderr, stdout bytes.Buffer
	cmd := exec.Command(command, args...)

	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	if stderr.String() != "" {
		return nil, errors.New(stderr.String())
	}

	return stdout.Bytes(), nil
}
