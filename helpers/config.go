package helpers

import (
	"embed"
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

	err = os.MkdirAll(filepath.Dir(configTargetPath), 0700)
	if err != nil {
		return err
	}

	_, err = os.Stat(configTargetPath)
	if err == nil {
		return nil
	} else if os.IsNotExist(err) {
		// The config holds api tokens, so it is a secrets file
		err := os.WriteFile(configTargetPath, content, 0600)
		if err != nil {
			return err
		}
		return nil
	} else {
		return err
	}
}
