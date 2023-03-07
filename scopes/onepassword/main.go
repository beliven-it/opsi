package onepassword

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func (o *onePassword) grantPermissions(vaultName string, userGroup string, permissions []string) error {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command("op", "vault", "group", "grant", "--vault", vaultName, "--group", userGroup, "--permissions", strings.Join(permissions, ","))
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

func (o *onePassword) listAccounts() ([]OnePasswordAccount, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command("op", "account", "list", "--format", "json")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, errors.New(stderr.String())
	}

	var listOfAccounts []OnePasswordAccount
	err = json.Unmarshal(stdout.Bytes(), &listOfAccounts)
	if err != nil {
		return nil, err
	}

	accountFilteredByAddress := []OnePasswordAccount{}
	for _, account := range listOfAccounts {
		if account.Url == o.address {
			accountFilteredByAddress = append(accountFilteredByAddress, account)
		}
	}

	return accountFilteredByAddress, nil
}

func (o *onePassword) listVaults() ([]OnePasswordVault, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command("op", "vault", "list", "--format", "json")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, errors.New(stderr.String())
	}

	var listOfVaults []OnePasswordVault
	err = json.Unmarshal(stdout.Bytes(), &listOfVaults)
	if err != nil {
		return nil, err
	}

	return listOfVaults, nil
}

func (o *onePassword) listGroups() ([]OnePasswordGroup, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command("op", "group", "list", "--format", "json")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, errors.New(stderr.String())
	}

	var listOfGroups []OnePasswordGroup
	err = json.Unmarshal(stdout.Bytes(), &listOfGroups)
	if err != nil {
		return nil, err
	}

	return listOfGroups, nil
}

func (o *onePassword) createGroup(groupName string) error {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("op", "group", "create", groupName)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

func (o *onePassword) createVault(vaultName string) error {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("op", "vault", "create", vaultName)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

func (o *onePassword) revokeUserFromGroup(groupName string) error {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("op", "group", "user", "revoke", "--user", o.account.UserUUID, "--group", groupName)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

func (o *onePassword) revokeUserFromVault(vaultName string) error {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("op", "vault", "user", "revoke", "--user", o.account.UserUUID, "--vault", vaultName)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

func (o *onePassword) addVault(vaultName string) error {
	listOfVaults, err := o.listVaults()
	if err != nil {
		return err
	}

	// Check if vault already exist
	for _, vault := range listOfVaults {
		if vault.Name == vaultName {
			return fmt.Errorf(
				"vault %s already exist",
				vaultName,
			)
		}
	}

	// Create vault
	err = o.createVault(vaultName)
	if err != nil {
		return err
	}

	return nil
}

func (o *onePassword) addGroup(groupName string) error {
	listOfGroups, err := o.listGroups()
	if err != nil {
		return err
	}

	// Check if group already exist
	for _, group := range listOfGroups {
		if group.Name == groupName {
			return fmt.Errorf(
				"group %s already exist",
				groupName,
			)
		}
	}

	// Create group
	err = o.createGroup(groupName)
	if err != nil {
		return err
	}

	return nil
}

func (o *onePassword) createContainer(projectName string, permissions []string) error {
	err := o.addGroup(projectName)
	if err != nil {
		return err
	}

	err = o.addVault(projectName)
	if err != nil {
		return err
	}

	err = o.grantPermissions(projectName, "Owners", privilegedPermissions)
	if err != nil {
		return err
	}

	err = o.grantPermissions(projectName, "Administrators", privilegedPermissions)
	if err != nil {
		return err
	}

	err = o.grantPermissions(projectName, projectName, permissions)
	if err != nil {
		return err
	}

	err = o.revokeUserFromGroup(projectName)
	if err != nil {
		return err
	}

	err = o.revokeUserFromVault(projectName)
	if err != nil {
		return err
	}

	return nil
}

func (o *onePassword) Create(projectName string) error {
	priName := projectName + " - PRI"
	pubName := projectName + " - PUB"

	listOfAccounts, err := o.listAccounts()
	if err != nil {
		return err
	}

	switch len(listOfAccounts) {
	case 1:
		o.account = listOfAccounts[0]
	case 0:
		fmt.Println("There's no account for this device, try to login and retry")
		return nil
	default:
		for {
			fmt.Println("Choose one of this accounts from the list below:")
			for index, account := range listOfAccounts {
				fmt.Printf("[%d] - %s", index, account.Email)
			}

			reader := bufio.NewReader(os.Stdin)

			value, _ := reader.ReadString('\n')
			valueAsInt, err := strconv.Atoi(value)
			if err != nil {
				return errors.New("invalid input")
			}

			if valueAsInt <= len(listOfAccounts)-1 {
				o.account = listOfAccounts[valueAsInt]
				break
			}
		}
	}

	// Create primary container
	err = o.createContainer(priName, unprivilegedPriPermissions)
	if err != nil {
		return err
	}

	// Create public container
	err = o.createContainer(pubName, unprivilegedPubPermissions)
	if err != nil {
		return err
	}

	// Grant permissions
	err = o.grantPermissions(pubName, priName, unprivilegedPriPermissions)
	if err != nil {
		return err
	}

	return nil
}

func (o *onePassword) Deprovisioning(userEmail string) error {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command("op", "user", "list", "--format", "json")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}

	var listOfUsers []OnePasswordUser
	err = json.Unmarshal(stdout.Bytes(), &listOfUsers)
	if err != nil {
		return err
	}

	userToSuspend := []OnePasswordUser{}
	for _, user := range listOfUsers {
		if userEmail != "" {
			if userEmail != user.Email {
				continue
			}
		}

		if user.State == "SUSPENDED" {
			userToSuspend = append(userToSuspend, user)
		}
	}

	if len(userToSuspend) == 0 {
		fmt.Println("There aren't any user to delete")
		return nil
	}

	fmt.Println("Are you sure to remove these users? (y/n)")
	for _, user := range userToSuspend {
		fmt.Printf("%s\n", user.Name)
	}

	var reader = bufio.NewReader(os.Stdin)
	confirmation, _ := reader.ReadString('\n')
	if confirmation != "y" {
		return nil
	}

	for _, user := range userToSuspend {
		fmt.Printf("%s\n", user.Name)
		// exec.Command("op", "user", "delete", user.ID)
	}

	return nil
}

func NewOnePassword(address string) OnePassword {
	return &onePassword{
		address: address,
	}
}
