package helpers

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func ExecuteScript(content []byte, args ...string) ([]byte, error) {
	// Create the temp script
	file, err := os.CreateTemp("/tmp", "script-")
	if err != nil {
		return nil, err
	}

	// defer os.Remove(file.Name())

	_, err = file.Write(content)
	if err != nil {
		return nil, err
	}

	os.Chmod(file.Name(), 0711)

	var stdout, stderr bytes.Buffer

	cmd := exec.Command(file.Name(), args...)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println(stderr.String())
		fmt.Println(stdout.String())

		fmt.Println(err)

		return nil, errors.New(stdout.String())
	}
	return stdout.Bytes(), nil

}
