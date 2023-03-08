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

	// If the list obtained is empty return the list of the projects collected until now.
	if len(list) == 0 {
		return entities, nil
	}

	// Otherwise continue to iterate the projects of the next page.
	entities = append(entities, list...)
	return g.walkThroughRequest(endpoint, entities, page+1)

}

func (g *gitlab) setupBranch(projectID int, branch gitlabSetupBranchRequest) error {
	// Create the endpoint
	endpoint := fmt.Sprintf("/projects/%d/protected_branches/%s", projectID, branch.Name)

	// Delete the branch for the specific project
	_, err := g.request("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	// Create the branch using the correct settings
	_, err = g.request("POST", fmt.Sprintf("/projects/%d/protected_branches", projectID), nil, map[string]string{
		"name":               branch.Name,
		"push_access_level":  strconv.Itoa(branch.PushAccessLevel),
		"merge_access_level": strconv.Itoa(branch.MergeAccessLevel),
	})
	if err != nil {
		return err
	}

	return nil
}

func (g *gitlab) setDefaultBranch(projectID int, branch string) error {
	var putPayload = map[string]interface{}{
		"default_branch":                branch,
		"ci_forward_deployment_enabled": false,
		"service_desk_enabled":          false,
	}

	_, err := g.request("PUT", fmt.Sprintf("/projects/%d", projectID), putPayload, nil)
	if err != nil {
		return err
	}

	return nil
}

// Apply a cleanUP policy on gitlab project.
func (g *gitlab) setCleanUpPolicy(projectID int) error {
	_, err := g.request("PUT", fmt.Sprintf("/projects/%d", projectID), defaultCleanUpPolicy, nil)

	return err
}

func (g *gitlab) CreateProject(name string, path string, groupID int) error {
	// Check if name and path
	// are property set
	if name == "" || path == "" {
		return errors.New("missing name or path arguments")
	}

	// Create the POST request payload
	payload := defaultGitlabCreatePayload
	payload.Name = name
	payload.Path = path
	payload.NamespaceID = groupID

	// Execute the request
	bodyResponse, err := g.request("POST", "/projects", payload, nil)
	if err != nil {
		return err
	}

	var project gitlabProjectResponse
	err = json.Unmarshal(bodyResponse, &project)
	if err != nil {
		return err
	}

	fmt.Println("The project ID is:", project.ID)

	// Create array of branches payload to create
	branches := []map[string]string{
		{
			"branch": "staging",
			"ref":    "master",
		},
		{
			"branch": "develop",
			"ref":    "staging",
		},
	}

	// Perform the request for create
	endpoint := fmt.Sprintf("/projects/%d/repository/branches", project.ID)
	for _, branch := range branches {
		_, err = g.request("POST", endpoint, nil, branch)
		if err != nil {
			return err
		}
	}

	// Setup the default branch of the project
	err = g.setDefaultBranch(project.ID, defaultBranch)
	if err != nil {
		return err
	}

	// Apply seggings to all the branches interested
	requestsPayload := []gitlabSetupBranchRequest{
		defaultProjectDevelopSettings,
		defaultProjectStagingSettings,
		defaultProjectMasterSettings,
	}

	for _, payload := range requestsPayload {
		err = g.setupBranch(project.ID, payload)
		if err != nil {
			return err
		}
	}

	err = g.setCleanUpPolicy(project.ID)
	if err != nil {
		return err
	}

	return nil
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

func (g *gitlab) DeleteEnvs(projectID string, env string, force bool) error {
	variables, err := g.listVariables(projectID, env)
	if err != nil {
		return err
	}

	if !force {
		if len(variables) == 0 {
			fmt.Println("No variables to delete")
			return nil
		}

		for _, variable := range variables {
			fmt.Printf("- %s [%s]\n", variable.Key, variable.EnvironmentScope)
		}

		rgx := regexp.MustCompile(`\n`)

		fmt.Println("The following variables will be deleted. Are you sure? (y/n)")
		reader := bufio.NewReader(os.Stdin)
		value, _ := reader.ReadString('\n')
		value = rgx.ReplaceAllString(value, "")
		if value != "y" {
			fmt.Println("Ok, abort procedure")
			return nil
		}
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
func (g *gitlab) CreateSubgroup(name string, path string, group *int) error {
	// Check if name and path are property set
	if name == "" || path == "" {
		return errors.New("missing name or path arguments")
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
		return err
	}

	// Take the response
	var subgroup gitlabSubgroupResponse
	err = json.Unmarshal(bodyResponse, &subgroup)
	if err != nil {
		return err
	}

	// Print ID of the group
	fmt.Println("The group ID is:", subgroup.ID)

	return nil
}

func (g *gitlab) BulkSettings() error {
	projects, err := g.walkThroughRequest("/projects", []gitlabEntityWithID{}, 0)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	messages := map[int][]string{}

	for _, project := range projects {
		wg.Add(1)
		func(projectID int) {
			defer wg.Done()

			messages[projectID] = []string{}

			messages[projectID] = append(messages[projectID], fmt.Sprintf("Update project #%d settings", projectID))

			response, err := g.request("GET", fmt.Sprintf("/projects/%d/repository/branches", projectID), nil, nil)
			if err != nil {
				messages[projectID] = append(
					messages[projectID],
					fmt.Sprintf("error fetching branches for project #%d", projectID),
				)
				return
			}

			var branches []gitlabBranchResponse

			err = json.Unmarshal(response, &branches)
			if err != nil {
				messages[projectID] = append(
					messages[projectID], fmt.Sprintf("cannot retrieve data about branches for project #%d", projectID),
				)
				return
			}

			var branchAsOctet = 0

			for _, branch := range branches {
				if branch.Name == "master" {
					branchAsOctet += 4
				}

				if branch.Name == "staging" {
					branchAsOctet += 2
				}

				if branch.Name == "develop" {
					branchAsOctet += 1
				}
			}

			var actions = []gitlabSetupBranchRequest{}

			if branchAsOctet >= 5 {
				actions = append(actions, gitlabSetupBranchRequest{
					Name:             "master",
					PushAccessLevel:  0,
					MergeAccessLevel: 40,
				})
			}

			if branchAsOctet == 7 || branchAsOctet == 3 {
				actions = append(actions, gitlabSetupBranchRequest{
					Name:             "staging",
					PushAccessLevel:  0,
					MergeAccessLevel: 30,
				})
			}

			if branchAsOctet == 6 {
				actions = append(actions, gitlabSetupBranchRequest{
					Name:             "staging",
					PushAccessLevel:  30,
					MergeAccessLevel: 30,
				})
			}

			if branchAsOctet == 4 {
				actions = append(actions, gitlabSetupBranchRequest{
					Name:             "master",
					PushAccessLevel:  40,
					MergeAccessLevel: 40,
				})
			}

			if branchAsOctet == 1 {
				actions = append(actions, gitlabSetupBranchRequest{
					Name:             "develop",
					PushAccessLevel:  30,
					MergeAccessLevel: 30,
				})
			}

			switch branchAsOctet {
			case 7, 5, 4:
				err = g.setDefaultBranch(projectID, "master")
				if err != nil {
					messages[projectID] = append(
						messages[projectID],
						fmt.Sprintf("Error on set default branch %s for project #%d: %s", "master", projectID, err.Error()),
					)
				}
			case 6, 3:
				err = g.setDefaultBranch(projectID, "staging")
				if err != nil {
					messages[projectID] = append(
						messages[projectID],
						fmt.Sprintf("Error on set default branch %s for project #%d: %s", "staging", projectID, err.Error()),
					)
				}
			case 1:
				err = g.setDefaultBranch(projectID, "develop")
				if err != nil {
					messages[projectID] = append(
						messages[projectID],
						fmt.Sprintf("Error on set default branch %s for project #%d: %s", "develop", projectID, err.Error()),
					)
				}

			}

		}(project.ID)
	}

	for _, msgs := range messages {
		for _, message := range msgs {
			helpers.Log(message)
		}
	}

	wg.Wait()
	return nil
}

// Handle deprovisioninig of a user
func (g *gitlab) Deprovionioning(username string) error {
	// Retrieve user ID by the username provided.
	bodyResponse, err := g.request("GET", "/users", nil, map[string]string{
		"username": username,
	})
	if err != nil {
		return err
	}

	// Convert response into a usable structure.
	var users []gitlabUser
	err = json.Unmarshal(bodyResponse, &users)
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

	// List all projects wher
	groups, err := g.walkThroughRequest("/groups", []gitlabEntityWithID{}, 0)
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

func NewGitlab(apiURL string, token string, groupID int) Gitlab {
	return &gitlab{
		apiURL:  apiURL,
		token:   token,
		groupID: groupID,
	}
}
