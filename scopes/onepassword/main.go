package onepassword

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"opsi/helpers"
	"os"
	"strconv"
	"strings"
)

// Execute any OP command.
func (o *onePassword) executeCommand(args ...string) ([]byte, error) {
	return execOp(args...)
}

// execOp runs the op binary. It is a variable so tests can replace it.
var execOp = func(args ...string) ([]byte, error) {
	return helpers.Exec("op", args...)
}

func (o *onePassword) grantPermissions(vaultName string, userGroup string, permissions []string) error {
	_, err := o.executeCommand("vault", "group", "grant", "--vault", vaultName, "--group", userGroup, "--permissions", strings.Join(permissions, ","))
	return err
}

func (o *onePassword) listAccounts() ([]OnePasswordAccount, error) {
	output, err := o.executeCommand("account", "list", "--format", "json")
	if err != nil {
		return nil, err
	}

	var listOfAccounts []OnePasswordAccount
	err = json.Unmarshal(output, &listOfAccounts)
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
	output, err := o.executeCommand("vault", "list", "--format", "json")
	if err != nil {
		return nil, err
	}

	var listOfVaults []OnePasswordVault
	err = json.Unmarshal(output, &listOfVaults)
	return listOfVaults, err
}

func (o *onePassword) listGroups() ([]OnePasswordGroup, error) {
	output, err := o.executeCommand("group", "list", "--format", "json")
	if err != nil {
		return nil, err
	}

	var listOfGroups []OnePasswordGroup
	err = json.Unmarshal(output, &listOfGroups)
	return listOfGroups, err
}

func (o *onePassword) createGroup(groupName string) error {
	_, err := o.executeCommand("group", "create", groupName)
	return err
}

func (o *onePassword) createVault(vaultName string) error {
	_, err := o.executeCommand("vault", "create", vaultName)
	return err
}

func (o *onePassword) revokeUserFromGroup(groupName string) error {
	_, err := o.executeCommand("group", "user", "revoke", "--user", o.account.UserUUID, "--group", groupName)
	return err
}

func (o *onePassword) revokeUserFromVault(vaultName string) error {
	_, err := o.executeCommand("vault", "user", "revoke", "--user", o.account.UserUUID, "--vault", vaultName)
	return err
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

	// Create vault and return output response
	return o.createVault(vaultName)
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
	return o.createGroup(groupName)
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

	return o.revokeUserFromVault(projectName)
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
	return o.grantPermissions(pubName, priName, unprivilegedPriPermissions)
}

// Deprovisioning deletes the suspended users from the workspace: all of them,
// or only the one matching userEmail when it is provided.
// With dryRun it only lists the users that would be deleted.
// The confirmation is up to the caller.
func (o *onePassword) Deprovisioning(userEmail string, dryRun bool) error {
	output, err := o.executeCommand("user", "list", "--format", "json")
	if err != nil {
		return err
	}

	var listOfUsers []OnePasswordUser
	err = json.Unmarshal(output, &listOfUsers)
	if err != nil {
		return err
	}

	userToDelete := []OnePasswordUser{}
	found := false
	for _, user := range listOfUsers {
		if userEmail != "" && !strings.EqualFold(userEmail, user.Email) {
			continue
		}
		found = true

		if user.State == "SUSPENDED" {
			userToDelete = append(userToDelete, user)
		}
	}

	if userEmail != "" && !found {
		return fmt.Errorf("user %q does not exist on 1password", userEmail)
	}

	if len(userToDelete) == 0 {
		if userEmail != "" {
			fmt.Printf("User %q is not suspended, nothing to delete\n", userEmail)
		} else {
			fmt.Println("There aren't any suspended users to delete")
		}
		return nil
	}

	if dryRun {
		fmt.Println("Dry run, these users would be deleted:")
		for _, user := range userToDelete {
			fmt.Printf("%s (%s)\n", user.Name, user.Email)
		}
		return nil
	}

	var failures []error
	deleted := 0
	for _, user := range userToDelete {
		if _, err := o.executeCommand("user", "delete", user.ID); err != nil {
			failures = append(failures, fmt.Errorf("%s (%s): %w", user.Name, user.Email, err))
			continue
		}

		deleted++
		fmt.Printf("Deleted %s (%s)\n", user.Name, user.Email)
	}

	fmt.Printf("Deleted %d of %d users\n", deleted, len(userToDelete))

	return errors.Join(failures...)
}

func NewOnePassword(address string) OnePassword {
	return &onePassword{
		address: address,
	}
}
