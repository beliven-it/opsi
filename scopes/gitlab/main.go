package gitlab

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"opsi/helpers"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// The request method perform an HTTP call into gitlab instance using
// the APIs endpoints.
func (g *gitlab) request(method string, endpoint string, body any, queryMap map[string]string) ([]byte, error) {
	return helpers.Request(method, g.apiURL+endpoint, body, queryMap, map[string]string{
		"Content-Type":  "application/json",
		"PRIVATE-TOKEN": g.token,
	})
}

func (g *gitlab) listUsers(filters map[string]string) ([]gitlabUser, error) {
	// Take the users
	response, err := g.request("GET", "/users", nil, filters)
	if err != nil {
		return nil, err
	}

	// Convert the output of the response in
	// an usable structure of data.
	var listOfUsers []gitlabUser
	err = json.Unmarshal(response, &listOfUsers)

	return listOfUsers, err
}

func (g *gitlab) listLeadUsers() ([]gitlabUser, error) {
	users, err := g.listUsers(nil)
	if err != nil {
		return nil, err
	}

	leadNoteRgx := regexp.MustCompile(gitlabDefaultGroupMember)
	leadUsers := []gitlabUser{}
	for _, user := range users {
		if leadNoteRgx.MatchString(user.Note) {
			leadUsers = append(leadUsers, user)
		}
	}

	return leadUsers, nil
}

func (g *gitlab) addUserToGroup(groupID int, userID int, accessLevel int) error {
	// Create the payload for the request
	payload := gitlabAddUserToGroupRequest{
		ID:          groupID,
		UserID:      userID,
		AccessLevel: accessLevel,
	}

	// Make the endpoint
	endpoint := fmt.Sprintf("/groups/%d/members", groupID)

	// Perform the request
	_, err := g.request("POST", endpoint, payload, nil)
	return err
}

// Take the list of variables for the specified project ID.
// Also, the output will be filtered for the environment provided.
func (g *gitlab) listVariables(projectID string, env string) ([]gitlabProjectListVariable, error) {
	response, err := g.request("GET", fmt.Sprintf("/projects/%s/variables", projectID), nil, nil)
	if err != nil {
		return nil, err
	}

	// Convert the output of the response in
	// an usable structure of data.
	var listOfVariables []gitlabProjectListVariable
	err = json.Unmarshal(response, &listOfVariables)
	if err != nil {
		return nil, err
	}

	// Check the environement.
	// If the environment variable provided
	// is all or a wildcard return the list without
	// filtering.
	if env == "all" || env == "*" {
		return listOfVariables, nil
	}

	// Otherwise filter the variables by the
	// environment provided.
	listOfVariablesFiltered := []gitlabProjectListVariable{}
	for _, variable := range listOfVariables {
		if variable.EnvironmentScope == env {
			listOfVariablesFiltered = append(listOfVariablesFiltered, variable)
		}
	}

	return listOfVariablesFiltered, nil
}

// This is a recursive function for move through the paginated endpoint of
// lists. The only way to perform this action is to iterate the endpoint
// until the list of project is empty.
func (g *gitlab) walkThroughRequest(endpoint string, entities []gitlabEntityWithID, page int) ([]gitlabEntityWithID, error) {
	// Convert the page numeric value to integer.
	// This is required because the query params structure accept
	// a map of string of strings.
	pageAsString := strconv.Itoa(page)

	// Perform the paginated http call.
	listAsBytes, err := g.request("GET", endpoint, nil, map[string]string{
		"page": pageAsString,
	})
	if err != nil {
		return entities, err
	}

	// Convert the response of the http call into a usable structure data.
	var list []gitlabEntityWithID
	err = json.Unmarshal(listAsBytes, &list)
	if err != nil {
		return entities, err
	}

	// If the list obtained is empty return the list of the items collected until now.
	if len(list) == 0 {
		return entities, nil
	}

	// Increase pagination
	page = page + 1

	// Otherwise continue to iterate the items of the next page.
	entities = append(entities, list...)
	return g.walkThroughRequest(endpoint, entities, page)
}

func (g *gitlab) reSetupBranch(projectID int, branch gitlabSetupBranchRequest) error {
	// Create the endpoint
	endpoint := fmt.Sprintf("/projects/%d/protected_branches/%s", projectID, branch.Name)

	// Delete the branch for the specific project
	g.request("DELETE", endpoint, nil, nil)

	return g.setupBranch(projectID, branch)
}

func (g *gitlab) setupBranch(projectID int, branch gitlabSetupBranchRequest) error {
	// Create the branch using the correct settings
	_, err := g.request("POST", fmt.Sprintf("/projects/%d/protected_branches", projectID), nil, map[string]string{
		"name":               branch.Name,
		"push_access_level":  strconv.Itoa(branch.PushAccessLevel),
		"merge_access_level": strconv.Itoa(branch.MergeAccessLevel),
	})

	return err
}

func (g *gitlab) createBranch(projectID int, branchPayload map[string]string) error {
	endpoint := fmt.Sprintf("/projects/%d/repository/branches", projectID)
	_, err := g.request("POST", endpoint, nil, branchPayload)
	return err
}

func (g *gitlab) setDefaultBranch(projectID int, branch string) error {
	var putPayload = map[string]interface{}{
		"default_branch":                branch,
		"ci_forward_deployment_enabled": false,
		"service_desk_enabled":          false,
	}

	_, err := g.request("PUT", fmt.Sprintf("/projects/%d", projectID), putPayload, nil)

	return err
}

// Apply a cleanUP policy on gitlab project.
func (g *gitlab) applyCleanUpPolicy(projectID int) error {
	_, err := g.request("PUT", fmt.Sprintf("/projects/%d", projectID), defaultCleanUpPolicy, nil)

	return err
}

func (g *gitlab) CreateProject(name string, path string, groupID int, defaultBranch string) (int, error) {
	// Create the POST request payload
	payload := defaultGitlabCreatePayload
	payload.Name = name
	payload.Path = path
	payload.NamespaceID = groupID

	// Execute the request
	bodyResponse, err := g.request("POST", "/projects", payload, nil)
	if err != nil {
		return 0, err
	}

	// Read the response
	var project gitlabProjectResponse
	err = json.Unmarshal(bodyResponse, &project)
	if err != nil {
		return 0, err
	}

	// Prepare an array of branches.
	// The branches contained in this array will be created
	// along the default branch.
	branches := []map[string]string{
		{
			"branch": "staging",
			"ref":    defaultBranch,
		},
		{
			"branch": "develop",
			"ref":    "staging",
		},
	}

	// Perform the request for create the branch
	for _, branch := range branches {
		err = g.createBranch(project.ID, branch)
		if err != nil {
			return 0, err
		}
	}

	// Set the default branch of the project
	err = g.setDefaultBranch(project.ID, defaultBranch)
	if err != nil {
		return 0, err
	}

	// Apply settings to all the branches interested
	requestsPayload := []gitlabSetupBranchRequest{
		defaultProjectDevelopSettings,
		defaultProjectStagingSettings,
		{
			Name:             defaultBranch,
			PushAccessLevel:  0,
			MergeAccessLevel: 40,
		},
	}

	for _, payload := range requestsPayload {
		err = g.setupBranch(project.ID, payload)
		if err != nil {
			return 0, err
		}
	}

	// Apply the cleanUP policy for the project created
	err = g.applyCleanUpPolicy(project.ID)
	if err != nil {
		return 0, err
	}

	return project.ID, nil
}

func (g *gitlab) CreateEnvs(projectID string, env string, envPath string) error {
	// Read the file env provided
	envFile, err := os.Open(envPath)
	if err != nil {
		return err
	}

	// Close the file once the function is finish
	defer envFile.Close()

	// Create a scanner for read the file line by line
	scanner := bufio.NewScanner(envFile)
	scanner.Split(bufio.ScanLines)

	// Prepare two regexp for check the existence of env vars with some pattern.
	// If the var start with MASKED, on Gitlab the var must be set as "masked"
	// If the var start with NOPROTECTED, on Gitlab the var must be set as "no protected"
	maskedRgx := regexp.MustCompile(`MASKED_`)
	unprotectedRgx := regexp.MustCompile(`NOPROTECTED_`)

	// Read line by line the buffer
	for scanner.Scan() {
		// Take line into string rapresentation
		text := scanner.Text()

		// Each variable have this format. KEY = VALUE
		// Separate the key from value splitting on character "="
		// if the matchs are more than 2 means the value contains one
		// or more "=" characters. Join these arguments again.
		partials := strings.Split(text, "=")
		if len(partials) < 2 {
			continue
		} else if len(partials) > 2 {
			partials[1] = strings.Join(partials[1:], "=")
		}

		// Take the KEY value
		key := partials[0]

		// Take the VAlUE
		value := partials[1]

		// Check if the key is masked
		isMasked := maskedRgx.MatchString(key)

		// Check if the key is unprotected
		isUnProtected := unprotectedRgx.MatchString(key)

		// Remove the unprotected and masked prefix from the key
		key = maskedRgx.ReplaceAllString(key, "")
		key = unprotectedRgx.ReplaceAllString(key, "")

		// Create the payload for the env creation
		payload := gitlabCreateEnvRequest{
			VariableType:     "env_var",
			Key:              key,
			Value:            value,
			Masked:           isMasked,
			Protected:        !isUnProtected,
			EnvironmentScope: env,
		}

		// Create the environment variable
		g.request("POST", fmt.Sprintf("/projects/%s/variables", projectID), payload, nil)
	}

	return nil
}

func (g *gitlab) ListEnvs(projectID string, env string) error {
	// Take the list of env
	variables, err := g.listVariables(projectID, env)
	if err != nil {
		return err
	}

	// Group the variables found by environment
	var groupedVariables = map[string][]gitlabProjectListVariable{}
	var maxKeyLength = 0
	for _, variable := range variables {
		_, ok := groupedVariables[variable.EnvironmentScope]
		if !ok {
			groupedVariables[variable.EnvironmentScope] = []gitlabProjectListVariable{}
		}

		if len(variable.Key) > maxKeyLength {
			maxKeyLength = len(variable.Key) + 1
		}

		groupedVariables[variable.EnvironmentScope] = append(groupedVariables[variable.EnvironmentScope], variable)
	}

	for key, variables := range groupedVariables {
		fmt.Printf("\n# [%s]\n", key)
		for _, variable := range variables {
			keySpacer := make([]string, maxKeyLength-len(variable.Key))

			fmt.Printf("%s%s = %s\n", variable.Key, strings.Join(keySpacer, " "), variable.Value)
		}
	}

	return nil
}

func (g *gitlab) DeleteEnvs(projectID string, env string) error {
	variables, err := g.listVariables(projectID, env)
	if err != nil {
		return err
	}

	for _, variable := range variables {
		queryParams := map[string]string{}
		queryParams["filter[environment_scope]"] = variable.EnvironmentScope

		_, err = g.request("DELETE", fmt.Sprintf("/projects/%s/variables/%s", projectID, variable.Key), nil, queryParams)
		if err != nil {
			return err
		}
	}

	return nil
}

// Create subgroup
func (g *gitlab) CreateSubgroup(name string, path string, group *int) (int, error) {
	// Check if name and path are property set
	if name == "" || path == "" {
		return 0, errors.New("missing name or path arguments")
	}

	// Create the POST request payload
	payload := gitlabCreateSubgroupRequest{
		Name:                  name,
		Path:                  path,
		ParentID:              group,
		Visibility:            "private",
		ProjectCreationLevel:  "maintainer",
		SubgroupCreationLevel: "owner",
		RequestAccessEnabled:  false,
	}

	// Execute the request
	bodyResponse, err := g.request("POST", "/groups", payload, nil)
	if err != nil {
		return 0, err
	}

	// Take the response
	var subgroup gitlabSubgroupResponse
	err = json.Unmarshal(bodyResponse, &subgroup)

	// If group is created at root level
	// Provide the lead users to the group
	if group == nil {
		leadUsers, err := g.listLeadUsers()
		if err != nil {
			return 0, err
		}

		for _, user := range leadUsers {
			g.addUserToGroup(subgroup.ID, user.ID, gitlabOwnerPermission)
		}
	}

	// Return the values
	return subgroup.ID, err
}

func (g *gitlab) BulkSettings(channel *chan string) error {
	projects, err := g.walkThroughRequest("/projects", []gitlabEntityWithID{}, 1)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	for _, project := range projects {
		wg.Add(1)
		func(projectID int) {
			defer wg.Done()

			*channel <- fmt.Sprintf("Update project #%d settings", projectID)

			response, err := g.request("GET", fmt.Sprintf("/projects/%d/repository/branches", projectID), nil, nil)
			if err != nil {
				*channel <- fmt.Sprintf("error fetching branches for project #%d", projectID)
				return
			}

			var branches []gitlabBranchResponse

			err = json.Unmarshal(response, &branches)
			if err != nil {
				*channel <- fmt.Sprintf("cannot retrieve data about branches for project #%d", projectID)
				return
			}

			var branchAsOctet = 0
			defaultBranch := ""
			for _, branch := range branches {
				if branch.Default {
					branchAsOctet += 4
					defaultBranch = branch.Name
				}

				if branch.Name == "staging" {
					branchAsOctet += 2
				}

				if branch.Name == "develop" {
					branchAsOctet += 1
				}
			}

			var actions = []gitlabSetupBranchRequest{}

			// Project has only default branch like master or main
			if branchAsOctet == 4 {
				actions = append(actions, gitlabSetupBranchRequest{
					Name:             defaultBranch,
					PushAccessLevel:  30,
					MergeAccessLevel: 30,
				})
			}

			// Project has other branch before the default
			if branchAsOctet > 4 {
				actions = append(actions, gitlabSetupBranchRequest{
					Name:             defaultBranch,
					PushAccessLevel:  0,
					MergeAccessLevel: 40,
				})
			}

			// Project has only staging branch
			if branchAsOctet == 2 {
				actions = append(actions, gitlabSetupBranchRequest{
					Name:             "staging",
					PushAccessLevel:  30,
					MergeAccessLevel: 30,
				})
			}

			// Project has other branch before the staging
			if branchAsOctet == 3 || branchAsOctet == 6 || branchAsOctet == 7 {
				actions = append(actions, gitlabSetupBranchRequest{
					Name:             "staging",
					PushAccessLevel:  0,
					MergeAccessLevel: 30,
				})
			}

			// Project has develop branch
			if branchAsOctet == 1 || branchAsOctet == 3 || branchAsOctet == 5 || branchAsOctet == 7 {
				actions = append(actions, gitlabSetupBranchRequest{
					Name:             "develop",
					PushAccessLevel:  30,
					MergeAccessLevel: 30,
				})
			}

			switch branchAsOctet {
			case 7, 5, 4:
				err = g.setDefaultBranch(projectID, defaultBranch)
				if err != nil {
					*channel <- fmt.Sprintf("Error on set default branch %s for project #%d: %s", defaultBranch, projectID, err.Error())
				} else {
					*channel <- fmt.Sprintf("Set default branch %s for project #%d", defaultBranch, projectID)
				}
			case 6, 3:
				err = g.setDefaultBranch(projectID, "staging")
				if err != nil {
					*channel <- fmt.Sprintf("Error on set default branch %s for project #%d: %s", "staging", projectID, err.Error())
				} else {
					*channel <- fmt.Sprintf("Set default branch %s for project #%d", "staging", projectID)
				}
			case 1:
				err = g.setDefaultBranch(projectID, "develop")
				if err != nil {
					*channel <- fmt.Sprintf("Error on set default branch %s for project #%d: %s", "develop", projectID, err.Error())
				} else {
					*channel <- fmt.Sprintf("Set default branch %s for project #%d", "develop", projectID)
				}
			}

			// Setup branches
			for _, action := range actions {
				err = g.reSetupBranch(projectID, action)
				if err != nil {
					*channel <- fmt.Sprintf("Error on setup branch for project #%d: %s", projectID, err.Error())
				}
			}

			// Apply cleanup policy
			err = g.applyCleanUpPolicy(projectID)
			if err != nil {
				*channel <- fmt.Sprintf("Error on apply cleanup policy for project #%d: %s", projectID, err.Error())
			}

		}(project.ID)
	}

	wg.Wait()
	return nil
}

// Handle deprovisioninig of a user
func (g *gitlab) Deprovionioning(username string) error {
	// Retrieve user ID by the username provided.
	users, err := g.listUsers(map[string]string{
		"username": username,
	})
	if err != nil {
		return err
	}

	// If the list is empty the script cannot continue.
	if len(users) == 0 {
		return errors.New("user not found")
	}

	// Take the ID of the first user found.
	// Tipically the result of the list must be one.
	userID := users[0].ID

	// List all projects
	groups, err := g.walkThroughRequest("/groups", []gitlabEntityWithID{}, 1)
	if err != nil {
		return err
	}
	// projects := g.walkThroughProjects([]gitlabEntityWithID{}, 0)

	wg := sync.WaitGroup{}

	// Remove from groups and subgroups
	for _, group := range groups {
		wg.Add(1)
		go func(groupID int) {
			defer wg.Done()

			g.request("DELETE", fmt.Sprintf("/groups/%d/members/%d", groupID, userID), nil, nil)

		}(group.ID)
	}

	// Remove from projects
	// for _, project := range projects {
	// 	wg.Add(1)
	// 	go func(projectID int) {
	// 		defer wg.Done()

	// 		g.request("DELETE", fmt.Sprintf("/projects/%d/members/%d", projectID, userID), nil, nil)

	// 	}(project.ID)
	// }

	wg.Wait()

	return nil
}

func NewGitlab(apiURL string, token string) Gitlab {
	return &gitlab{
		apiURL: apiURL,
		token:  token,
	}
}
