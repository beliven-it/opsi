package helpers

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

func ConfigInit(template embed.FS, path string) error {
	content, err := template.ReadFile("config/template.yml")
	if err != nil {
		return err
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configTargetPath := userHomeDir + path

	err = os.MkdirAll(filepath.Dir(configTargetPath), 0755)
	if err != nil {
		return err
	}

	_, err = os.Stat(configTargetPath)
	if err == nil {
		return nil
	} else if os.IsNotExist(err) {
		err := os.WriteFile(configTargetPath, content, 0755)
		if err != nil {
			return err
		}
		return nil
	} else {
		fmt.Println("AAA", err.Error())
		return err
	}
}
