package scopes

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

type OnePassword struct {
	address string
	account OnePasswordAccount
}

type OnePasswordUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Type  string `json:"type"`
	State string `json:"state"`
}

type OnePasswordAccount struct {
	Url         string `json:"url"`
	Email       string `json:"email"`
	UserUUID    string `json:"user_uuid"`
	AccountUUID string `json:"account_uuid"`
}

type OnePasswordGroup struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	State     string `json:"state"`
	CreatedAt string `json:"created_at"`
}

type OnePasswordVault struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ContentVersion string `json:"content_version"`
}

var privilegedPermissions = []string{
	"view_items",
	"create_items",
	"edit_items",
	"archive_items",
	"delete_items",
	"view_and_copy_passwords",
	"view_item_history",
	"import_items",
	"export_items",
	"copy_and_share_items",
	"print_items",
	"manage_vault",
}

var unprivilegedPriPermissions = []string{
	"view_items",
	"create_items",
	"edit_items",
	"archive_items",
	"view_and_copy_passwords",
	"view_item_history",
	"import_items",
	"export_items",
	"copy_and_share_items",
}

var unprivilegedPubPermissions = []string{
	"view_items",
	"view_and_copy_passwords",
	"view_item_history",
	"copy_and_share_items",
}

func (o *OnePassword) grantPermissions(vaultName string, userGroup string, permissions []string) error {
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

func (o *OnePassword) listAccounts() ([]OnePasswordAccount, error) {
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

func (o *OnePassword) listVaults() ([]OnePasswordVault, error) {
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

func (o *OnePassword) listGroups() ([]OnePasswordGroup, error) {
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

func (o *OnePassword) createGroup(groupName string) error {
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

func (o *OnePassword) createVault(vaultName string) error {
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

func (o *OnePassword) revokeUserFromGroup(groupName string) error {
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

func (o *OnePassword) revokeUserFromVault(vaultName string) error {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("op", "vault", "user", "revoke", "--user", o.account.UserUUID, "--group", vaultName)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

func (o *OnePassword) addVault(vaultName string) error {
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

	// Revoke user from created group
	err = o.revokeUserFromVault(vaultName)
	if err != nil {
		return err
	}

	return nil
}

func (o *OnePassword) addGroup(groupName string) error {
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

	// Revoke user from created group
	err = o.revokeUserFromGroup(groupName)
	if err != nil {
		return err
	}

	return nil
}

func (o *OnePassword) createContainer(projectName string, permissions []string) error {
	// Create the primary group
	err := o.addGroup(projectName)
	if err != nil {
		return err
	}

	// Create the primary vault
	err = o.addVault(projectName)
	if err != nil {
		return err
	}

	// Grant vault permissions
	err = o.grantPermissions(projectName, "Owner", privilegedPermissions)
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

	return nil
}

func (o *OnePassword) Create(projectName string) error {
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
	err = o.grantPermissions(pubName, priName, privilegedPermissions)
	if err != nil {
		return err
	}

	return nil
}

func (o *OnePassword) Deprovisioning(userEmail string) error {
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
	return OnePassword{
		address: address,
	}
}
