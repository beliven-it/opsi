package helpers

import (
	"os"
	"os/exec"
)

func ExecuteScript(content []byte, args ...string) ([]byte, error) {
	// Create the temp script
	file, err := os.CreateTemp("/tmp", "script-")
	if err != nil {
		return nil, err
	}

	defer os.Remove(file.Name())

	_, err = file.Write(content)
	if err != nil {
		return nil, err
	}

	os.Chmod(file.Name(), 0711)
	return exec.Command(file.Name(), args...).Output()

}
