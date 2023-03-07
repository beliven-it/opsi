package hosts

type hostHSSH struct {
	Name     string
	Hostname string
	User     string
	Port     int
	Identity string
}

type host struct{}

type Host interface {
	CheckReboot() error
}
