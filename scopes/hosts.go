package scopes

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Hosts struct {
}

func (o *Hosts) createConnectionString(host hostHSSH) string {
	return fmt.Sprintf("%s@%s", host.User, host.Hostname)
}

func (o *Hosts) createPortString(host hostHSSH) string {
	if host.Port == 0 {
		return "22"
	}

	return fmt.Sprintf("%d", host.Port)
}

func (o *Hosts) listHosts() ([]hostHSSH, []string, error) {

	var stdout, stderr bytes.Buffer

	cmd := exec.Command("hssh", "l")
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return nil, nil, err
	}

	reader := bytes.NewReader(stdout.Bytes())
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	list := []hostHSSH{}
	ignoredHosts := []string{}

	for scanner.Scan() {
		text := scanner.Text()
		partials := strings.Split(text, "->")
		if len(partials) == 0 {
			continue
		}

		hostname := strings.Trim(partials[0], " ")

		host, err := o.findHost(hostname)
		if err != nil {
			ignoredHosts = append(ignoredHosts, hostname)
			continue
		}

		list = append(list, host)
	}

	return list, ignoredHosts, nil
}

func (o *Hosts) findHost(hostname string) (hostHSSH, error) {
	host := hostHSSH{}
	var stdout, stderr bytes.Buffer

	cmd := exec.Command("hssh", "f", hostname)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()

	if err != nil {
		return host, err
	}

	reader := bytes.NewReader(stdout.Bytes())
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		partials := strings.Split(line, ":")

		value := ""
		if len(partials) >= 2 {
			value = strings.Trim(partials[1], " ")
		}

		key := strings.Trim(partials[0], " ")

		switch key {
		case "Hostname":
			host.Hostname = value
		case "Port":
			host.Port, _ = strconv.Atoi(value)
		case "User":
			host.User = value
		case "Identity":
			host.Identity = value
		case "Name":
			host.Name = value
		}
	}

	return host, nil
}

func (o *Hosts) CheckReboot() error {
	listErrors := []string{}
	listRebootable := []string{}
	listUnrebootable := []string{}

	hosts, listIgnored, err := o.listHosts()
	if err != nil {
		return err
	}

	statusRgx := regexp.MustCompile(`exit status ([0-9]+)`)

	for _, host := range hosts {
		var stdout, stderr bytes.Buffer
		command := []string{
			o.createConnectionString(host),
			"-p",
			o.createPortString(host),
			"-o",
			"StrictHostKeyChecking=no",
			"-o",
			"PasswordAuthentication=no",
			"-o",
			"ConnectTimeout=5",
			`ls -la /var/run/reboot-required`,
		}

		cmd := exec.Command("ssh", command...)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err = cmd.Run()

		if err != nil {
			statusCode := statusRgx.ReplaceAllString(err.Error(), "$1")

			if statusCode == "2" || statusCode == "1" {
				listUnrebootable = append(listUnrebootable, host.Name)
			} else {
				listErrors = append(listErrors, host.Name)
			}

		} else {
			listRebootable = append(listRebootable, host.Name)
		}
	}

	fmt.Println("REBOOTABLE")
	for _, host := range listRebootable {
		fmt.Println("-", host)
	}

	fmt.Println("\nUN-REBOOTABLE")
	for _, host := range listUnrebootable {
		fmt.Println("-", host)
	}

	fmt.Println("\nSSH ERRORS")
	for _, host := range listErrors {
		fmt.Println("-", host)
	}

	fmt.Println("\nIGNORED FOR OTHER ERRORS")
	for _, host := range listIgnored {
		fmt.Println("-", host)
	}

	return nil

}

func NewHosts() Hosts {
	return Hosts{}
}
