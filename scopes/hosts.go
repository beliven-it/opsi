package scopes

import "opsi/helpers"

type Hosts struct {
}

func (o *Hosts) CheckReboot(scriptContent []byte, args ...string) ([]byte, error) {
	return helpers.ExecuteScript(scriptContent, args...)
}

func NewHosts() Hosts {
	return Hosts{}
}
