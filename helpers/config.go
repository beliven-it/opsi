package helpers

import (
	"embed"
	"os"
)

func ConfigInit(template embed.FS) error {
	content, err := template.ReadFile("config/template.yml")
	if err != nil {
		return err
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configTargetPath := userHomeDir + "/.opsi.yml"

	_, err = os.Stat(configTargetPath)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(configTargetPath)
			if err != nil {
				return err
			}

			_, err = file.Write(content)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
