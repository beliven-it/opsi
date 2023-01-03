package scopes

import (
	"opsi/helpers"
)

type OnePassword struct {
	address string
}

func (o *OnePassword) Create(scriptContent []byte, args ...string) ([]byte, error) {
	args = append(args, o.address)
	return helpers.ExecuteScript(scriptContent, args...)
}

func NewOnePassword(address string) OnePassword {
	return OnePassword{
		address: address,
	}
}
