package scopes

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

type Gitlab struct {
	token   string
	groupID int
	apiURL  string
}

func (g *Gitlab) request(method string, endpoint string, body []byte, queryMap map[string]string) ([]byte, error) {
	return helpers.Request(method, g.apiURL+endpoint, body, queryMap, map[string]string{
		"Content-Type":  "application/json",
		"PRIVATE-TOKEN": g.token,
	})
}

func (g *Gitlab) listVariables(projectID string, env string) ([]gitlabProjectListVariable, error) {
	response, err := g.request("GET", fmt.Sprintf("/projects/%s/variables", projectID), nil, nil)
	if err != nil {
		return nil, err
	}

	var listOfVariables []gitlabProjectListVariable

	err = json.Unmarshal(response, &listOfVariables)
	if err != nil {
		return nil, err
	}

	if env == "all" || env == "*" {
		return listOfVariables, nil
	}

	listOfVariablesFiltered := []gitlabProjectListVariable{}
	for _, variable := range listOfVariables {
		if variable.EnvironmentScope == env {
			listOfVariablesFiltered = append(listOfVariablesFiltered, variable)
		}
	}

	return listOfVariablesFiltered, nil
}

func (g *Gitlab) walkThroughProjects(projects []gitlabEntityWithID, page int) []gitlabEntityWithID {
	pageAsString := strconv.Itoa(page)
	listAsBytes, err := g.request("GET", "/projects", nil, map[string]string{
		"page": pageAsString,
	})

	if err != nil {
		return projects
	}

	var list []gitlabEntityWithID
	err = json.Unmarshal(listAsBytes, &list)

	if err != nil {
		return projects
	}

	if len(list) == 0 {
		return projects
	} else {
		projects = append(projects, list...)

		return g.walkThroughProjects(projects, page+1)
	}
}

func (g *Gitlab) walkThroughGroups(groups []gitlabEntityWithID, page int) []gitlabEntityWithID {
	pageAsString := strconv.Itoa(page)
	listAsBytes, err := g.request("GET", "/groups", nil, map[string]string{
		"page": pageAsString,
	})

	if err != nil {
		return groups
	}

	var list []gitlabEntityWithID
	err = json.Unmarshal(listAsBytes, &list)
	if err != nil {
		return groups
	}

	if len(list) == 0 {
		return groups
	} else {
		groups = append(groups, list...)

		return g.walkThroughGroups(groups, page+1)
	}
}

func (g *Gitlab) setupBranch(projectID int, branch gitlabSetupBranchRequest) error {

	endpoint := fmt.Sprintf("/projects/%d/protected_branches/%s", projectID, branch.Name)

	g.request("DELETE", endpoint, nil, nil)
	_, err := g.request("POST", fmt.Sprintf("/projects/%d/protected_branches", projectID), nil, map[string]string{
		"name":               branch.Name,
		"push_access_level":  strconv.Itoa(branch.PushAccessLevel),
		"merge_access_level": strconv.Itoa(branch.MergeAccessLevel),
	})
	if err != nil {
		return err
	}
	return nil
}

func (g *Gitlab) setDefaultBranch(projectID int, branch string) error {
	var putPayload = map[string]interface{}{
		"default_branch":                branch,
		"ci_forward_deployment_enabled": false,
		"service_desk_enabled":          false,
	}
	payloadAsBytes, err := json.Marshal(putPayload)
	if err != nil {
		return err
	}

	_, err = g.request("PUT", fmt.Sprintf("/projects/%d", projectID), payloadAsBytes, nil)
	if err != nil {
		return err
	}

	return nil
}

func (g *Gitlab) CreateProject(name string, path string, subgroupID int) error {
	// Check if name and path
	// are property set
	if name == "" || path == "" {
		return errors.New("missing name or path arguments")
	}

	// Create the POST request payload
	payload := gitlabProjectRequest{
		Name:                         name,
		Path:                         path,
		NamespaceID:                  subgroupID,
		MergeMethod:                  "ff",
		AnalyticsAccessLevel:         "disabled",
		SecurityAndComplianceEnabled: false,
		IssuesEnabled:                false,
		ForkingAccessLevel:           "disabled",
		LFSEnabled:                   false,
		WikiEnabled:                  false,
		PagesAccessLevel:             "disabled",
		OperationsAccessLevel:        "disabled",
		SharedRunnersEnabled:         false,
		InitializeWithReadME:         true,
		SquashOption:                 "never",
	}

	payloadAsBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Execute the request
	bodyResponse, err := g.request("POST", "/projects", payloadAsBytes, nil)
	if err != nil {
		return err
	}

	var project gitlabProjectResponse
	err = json.Unmarshal(bodyResponse, &project)
	if err != nil {
		return err
	}

	fmt.Println("The project ID is:", project.ID)

	_, err = g.request("POST", fmt.Sprintf("/projects/%d/repository/branches", project.ID), nil, map[string]string{
		"branch": "staging",
		"ref":    "master",
	})
	if err != nil {
		return err
	}

	_, err = g.request("POST", fmt.Sprintf("/projects/%d/repository/branches", project.ID), nil, map[string]string{
		"branch": "develop",
		"ref":    "staging",
	})
	if err != nil {
		return err
	}

	err = g.setDefaultBranch(project.ID, "master")
	if err != nil {
		return err
	}

	err = g.setupBranch(project.ID, gitlabSetupBranchRequest{
		Name:             "develop",
		PushAccessLevel:  30,
		MergeAccessLevel: 30,
	})
	if err != nil {
		return err
	}

	err = g.setupBranch(project.ID, gitlabSetupBranchRequest{
		Name:             "staging",
		PushAccessLevel:  0,
		MergeAccessLevel: 30,
	})
	if err != nil {
		return err
	}

	err = g.setupBranch(project.ID, gitlabSetupBranchRequest{
		Name:             "master",
		PushAccessLevel:  0,
		MergeAccessLevel: 30,
	})
	if err != nil {
		return err
	}

	return nil
}

func (g *Gitlab) CreateEnvs(projectID string, env string, envPath string) error {
	// Read file env provided
	envFile, err := os.Open(envPath)
	if err != nil {
		return err
	}

	defer envFile.Close()

	scanner := bufio.NewScanner(envFile)
	scanner.Split(bufio.ScanLines)

	maskedRgx := regexp.MustCompile(`MASKED_`)
	unprotectedRgx := regexp.MustCompile(`NOPROTECTED_`)

	for scanner.Scan() {
		text := scanner.Text()
		partials := strings.Split(text, "=")
		if len(partials) < 2 {
			continue
		} else if len(partials) > 2 {
			partials[1] = strings.Join(partials[1:], "=")
		}

		key := partials[0]
		value := partials[1]

		isMasked := maskedRgx.MatchString(key)
		isUnProtected := unprotectedRgx.MatchString(key)

		// Remove prefixes if any
		key = maskedRgx.ReplaceAllString(key, "")
		key = unprotectedRgx.ReplaceAllString(key, "")

		payload := gitlabEnvRequest{
			VariableType:     "env_var",
			Key:              key,
			Value:            value,
			Masked:           isMasked,
			Protected:        !isUnProtected,
			EnvironmentScope: env,
		}

		payloadAsBytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		g.request("POST", fmt.Sprintf("/projects/%s/variables", projectID), payloadAsBytes, nil)
	}

	return nil
}

func (g *Gitlab) ListEnvs(projectID string, env string) error {
	variables, err := g.listVariables(projectID, env)
	if err != nil {
		return err
	}

	if len(variables) == 0 {
		fmt.Println("No variables to delete")
		return nil
	}

	// Group env by environment
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

func (g *Gitlab) DeleteEnvs(projectID string, env string, force bool) error {
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
func (g *Gitlab) CreateSubgroup(name string, path string, group int) error {
	// Check if name and path
	// are property set
	if name == "" || path == "" {
		return errors.New("missing name or path arguments")
	}

	// Create the POST request payload
	payload := gitlabSubgroupRequest{
		Name:                  name,
		Path:                  path,
		ParentID:              group,
		Visibility:            "private",
		ProjectCreationLevel:  "maintainer",
		SubgroupCreationLevel: "owner",
		RequestAccessEnabled:  false,
	}

	payloadAsBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Execute the request
	bodyResponse, err := g.request("POST", "/groups", payloadAsBytes, nil)
	if err != nil {
		return err
	}

	var subgroup gitlabSubgroupResponse
	err = json.Unmarshal(bodyResponse, &subgroup)
	if err != nil {
		return err
	}

	fmt.Println("The subgroup ID is:", subgroup.ID)

	return nil
}

func (g *Gitlab) BulkSettings() error {
	projects := g.walkThroughProjects([]gitlabEntityWithID{}, 0)
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

			if branchAsOctet == 7 && branchAsOctet == 3 {
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
func (g *Gitlab) Deprovionioning(username string) error {
	// Retrieve user ID by username
	bodyResponse, err := g.request("GET", "/users", nil, map[string]string{
		"username": username,
	})
	if err != nil {
		return err
	}

	var users []gitlabUser
	err = json.Unmarshal(bodyResponse, &users)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return errors.New("user not found")
	}

	userID := users[0].ID

	// List all projects wher
	groups := g.walkThroughGroups([]gitlabEntityWithID{}, 0)
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
	return Gitlab{
		apiURL:  apiURL,
		token:   token,
		groupID: groupID,
	}
}
