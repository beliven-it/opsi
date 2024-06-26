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
	"time"
)

// The request method perform an HTTP call into gitlab instance using
// the APIs endpoints.
func (g *gitlab) request(method string, endpoint string, body any, queryMap map[string]string) ([]byte, error) {
	return helpers.Request(method, g.apiURL+endpoint, body, queryMap, map[string]string{
		"Content-Type":  "application/json",
		"PRIVATE-TOKEN": g.token,
	})
}

func (g *gitlab) mirrorRequest(method string, endpoint string, body any, queryMap map[string]string) ([]byte, error) {
	return helpers.Request(method, g.mirror.ApiURL+endpoint, body, queryMap, map[string]string{
		"Content-Type":  "application/json",
		"PRIVATE-TOKEN": g.mirror.Token,
	})
}

func (g *gitlab) viewGroup(groupID int) (gitlabSubgroupResponse, error) {
	var data gitlabSubgroupResponse
	endpoint := fmt.Sprintf("/groups/%d", groupID)
	response, err := g.request("GET", endpoint, nil, nil)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(response, &data)

	return data, err
}

func (g *gitlab) enableMirrorForProject(projectID int, projectName string) ([]byte, error) {
	endpoint := fmt.Sprintf("/projects/%d/remote_mirrors", projectID)

	payload := gitlabCreateMirrorRequest{
		Enabled:               true,
		OnlyProtectedBranched: true,
		URL:                   fmt.Sprintf("https://%s:%s@%s/%s.git", g.mirror.Username, g.mirror.Token, g.mirror.GroupPath, projectName),
	}

	return g.request("POST", endpoint, payload, nil)
}

func (g *gitlab) UpdateMirroring() error {
	//Retrieve projects list
	projectsList, err := g.listProjects()

	if err != nil {
		return err
	}

	var projectsIDWithMirroring []int
	var mirroringProjects []gitlabMirrorResponse

	// Filter repositories with mirroring enabled
	for _, project := range projectsList {
		mirroringProject, hasMirroring, err := g.checkMirroringExistence(project.ID)
		if err != nil {
			return err
		}
		if hasMirroring {
			projectsIDWithMirroring = append(projectsIDWithMirroring, project.ID)
			mirroringProjects = append(mirroringProjects, mirroringProject)
		}
	}

	//Recreate mirroring
	for i := 0; i < len(projectsIDWithMirroring) && i < len(mirroringProjects); i++ {
		//Delete current mirroring
		projectID := projectsIDWithMirroring[i]
		mirroringProjectID := mirroringProjects[i].ID
		err := g.deleteMirroring(projectID, mirroringProjectID)

		if err != nil {
			return err
		}
		// //Create new mirroring
		patternUrl := `\/([^\/]+)\.git$`
		re := regexp.MustCompile(patternUrl)
		matches := re.FindStringSubmatch(mirroringProjects[i].Url)
		projectName := matches[1]

		_, err = g.enableMirrorForProject(projectID, projectName)
		if err != nil {
			fmt.Printf("Error when updating mirroring for %s: %v\n", projectName, err)
		}
		fmt.Println("Mirroring updated for", projectName)
	}

	// Return the error
	return err
}

func (g *gitlab) listProjects() ([]gitlabProjectResponse, error) {
	fmt.Println("Retrieving projects list...")
	var allProjects []gitlabProjectResponse
	nextPage := 1
	perPage := 100

	for {
		// Construct the URL with pagination parameters
		endpoint := fmt.Sprintf("/projects?per_page=%d&page=%d&simple=true", perPage, nextPage)

		// Make the request
		response, err := g.request("GET", endpoint, nil, nil)
		if err != nil {
			return nil, err
		}

		// Parse the response and append projects to the list
		var projects []gitlabProjectResponse
		if err := json.Unmarshal(response, &projects); err != nil {
			return nil, err
		}
		allProjects = append(allProjects, projects...)

		// Check if there are more pages
		nextPage++
		if len(projects) < perPage {
			break // No more pages
		}
	}

	return allProjects, nil
}

func (g *gitlab) checkMirroringExistence(projectID int) (gitlabMirrorResponse, bool, error) {
	endpoint := fmt.Sprintf("/projects/%d/remote_mirrors", projectID)
	response, err := g.request("GET", endpoint, nil, nil)

	if err != nil {
		return gitlabMirrorResponse{}, false, err
	}

	// Parse the response to check if mirroring is enabled
	// Assuming the response contains JSON with a field indicating mirroring status
	var projectRemoteMirrors []gitlabMirrorResponse
	err = json.Unmarshal(response, &projectRemoteMirrors)
	if err != nil {
		return gitlabMirrorResponse{}, false, err
	}

	// Check if the project has at least one mirroring enabled
	for _, projectRemoteMirror := range projectRemoteMirrors {
		if projectRemoteMirror.Enabled {
			return projectRemoteMirror, true, nil
		}
	}
	// If no mirror response has mirroring enabled, return false
	return gitlabMirrorResponse{}, false, nil
}

func (g *gitlab) deleteMirroring(projectID int, mirroringProjectID int) error {
	endpoint := fmt.Sprintf("/projects/%d/remote_mirrors/%d", projectID, mirroringProjectID)
	_, err := g.request("DELETE", endpoint, nil, nil)

	return err
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

func (g *gitlab) listDefaultUsers(tipology string) ([]gitlabUser, error) {
	users, err := g.listUsers(map[string]string{"per_page": "100"})
	if err != nil {
		return nil, err
	}

	leadUsers := []gitlabUser{}
	for _, user := range users {
		if strings.EqualFold(tipology, user.Note) {
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

func (g *gitlab) UpdateCleanUpPolicy(projectID string) error {
	var err error
	if projectID != "" {
		id, err := strconv.Atoi(projectID)

		if err != nil {
			return err
		}

		err = g.applyCleanUpPolicy(id)
		if err != nil {
			return err
		}
	} else {
		projectsList, err := g.listProjects()

		if err != nil {
			return err
		}
		for _, project := range projectsList {
			err := g.applyCleanUpPolicy(project.ID)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// Function to check if an array contains an element
func (g *gitlab) contains(array []int, item int) (bool, error) {
	for i := 0; i < len(array); i++ {
		// check
		if array[i] == item {
			return true, nil
		}
	}
	return false, nil
}

// Apply a cleanUP policy on gitlab project.
func (g *gitlab) applyCleanUpPolicy(projectID int) error {
	isExcluded, err := g.contains(g.exclusions.CleanupPolicies, projectID)
	if !isExcluded {
		_, err = g.request("PUT", fmt.Sprintf("/projects/%d", projectID), defaultCleanUpPolicy, nil)
		fmt.Println("Cleanup policy updated for the project with ID", projectID)
	} else {
		fmt.Println("Cleanup policy not updated for the project with ID", projectID)
	}

	return err
}

// Set protected tags
func (g *gitlab) setupTag(projectID int) error {
	_, err := g.request("POST", fmt.Sprintf("/projects/%d/protected_tags", projectID), defaultProtectedTags, nil)

	return err
}

func (g *gitlab) createProject(options ProjectRequest) (gitlabProjectResponse, error) {
	var project gitlabProjectResponse

	payload := defaultGitlabCreatePayload
	payload.Visibility = options.Visibility
	payload.Name = options.Name
	payload.Path = options.Path
	payload.NamespaceID = options.Group
	payload.SharedRunnersEnabled = options.SharedRunners

	bodyResponse, err := g.request("POST", projectEndpoint, payload, nil)
	if err != nil {
		return project, err
	}

	// Read the response
	err = json.Unmarshal(bodyResponse, &project)
	return project, err
}

func (g *gitlab) createMirrorProject(options ProjectRequest) (gitlabProjectResponse, error) {
	var project gitlabProjectResponse

	payload := defaultGitlabMirrorCreatePayload
	payload.Name = options.Name
	payload.Path = options.Path
	payload.NamespaceID = options.Group

	bodyResponse, err := g.mirrorRequest("POST", projectEndpoint, payload, nil)
	if err != nil {
		return project, err
	}

	// Read the response
	err = json.Unmarshal(bodyResponse, &project)
	return project, err
}

func (g *gitlab) setupMirrorProject(name string, path string, groupID int) error {
	mirrorRequest := ProjectRequest{
		Name:  name,
		Path:  path,
		Group: groupID,
	}

	mirrorProject, err := g.createMirrorProject(mirrorRequest)
	if err != nil {
		return err
	}

	// Sleep because sometimes the repo seems not completed yet.
	// and the next call explode!!
	time.Sleep(2 * time.Second)

	endpoint := fmt.Sprintf("/projects/%d/protected_branches/main", mirrorProject.ID)
	payload := map[string]interface{}{
		"allow_force_push": true,
	}

	_, err = g.mirrorRequest("PATCH", endpoint, payload, nil)
	return err
}

func (g *gitlab) CreateProject(options ProjectRequest) (int, error) {
	// Take the group informations
	groupDetail, err := g.viewGroup(options.Group)
	if err != nil {
		return 0, err
	}

	// Inherit some attributes from the group
	if options.Visibility == "" {
		options.Visibility = groupDetail.Visibility
	}

	// Create the project
	projectRequest := options
	project, err := g.createProject(projectRequest)
	if err != nil {
		return 0, err
	}

	// Prepare an array of branches.
	// The branches contained in this array will be created
	// along the default branch.
	branches := []map[string]string{
		{
			"branch": "staging",
			"ref":    options.DefaultBranch,
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
	err = g.setDefaultBranch(project.ID, options.DefaultBranch)
	if err != nil {
		return 0, err
	}

	// Apply settings to all the branches interested
	requestsPayload := []gitlabSetupBranchRequest{
		defaultProjectDevelopSettings,
		defaultProjectStagingSettings,
		{
			Name:             options.DefaultBranch,
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

	// Set protected branches for tags
	err = g.setupTag(project.ID)
	if err != nil {
		return 0, err
	}

	// Apply the cleanUP policy for the project created
	err = g.applyCleanUpPolicy(project.ID)
	if err != nil {
		return 0, err
	}

	// If mirror is enables create the mirror project
	if options.Mirror {
		err = g.setupMirrorProject(options.Name, options.Path, g.mirror.GroupID)
		if err != nil {
			return 0, err
		}

		// Setup mirror project
		g.enableMirrorForProject(project.ID, options.Path)
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
func (g *gitlab) createGroup(payload gitlabCreateSubgroupRequest) (int, error) {
	// Check if name and path are property set
	if payload.Name == "" || payload.Path == "" {
		return 0, errors.New("missing name or path arguments")
	}

	// Execute the request
	bodyResponse, err := g.request("POST", "/groups", payload, nil)
	if err != nil {
		return 0, err
	}

	// Take the response
	var subgroup gitlabSubgroupResponse
	err = json.Unmarshal(bodyResponse, &subgroup)

	// If the group is create at root level,
	// provide some default users to the group
	typeOfUsers := []gitlabDefaultUser{
		{tipology: gitlabDefaultGroupMemberDeveloper, permission: gitlabDeveloperPermission},
		{tipology: gitlabDefaultGroupMemberMaintainer, permission: gitlabMaintainerPermission},
		{tipology: gitlabDefaultGroupMemberOwner, permission: gitlabOwnerPermission},
		{tipology: gitlabDefaultGroupMember, permission: gitlabOwnerPermission},
	}

	for _, t := range typeOfUsers {
		defaultUsersOfType, err := g.listDefaultUsers(t.tipology)
		if err != nil {
			return 0, err
		}

		for _, user := range defaultUsersOfType {
			g.addUserToGroup(subgroup.ID, user.ID, t.permission)
		}
	}

	// Return the values
	return subgroup.ID, err
}

// Create group
func (g *gitlab) CreateGroup(name string, path string, visibility string) (int, error) {
	payload := gitlabCreateSubgroupRequest{
		Name:                  name,
		Path:                  path,
		ParentID:              nil,
		Visibility:            visibility,
		RequestAccessEnabled:  false,
		ProjectCreationLevel:  "maintainer",
		SubgroupCreationLevel: "owner",
	}

	return g.createGroup(payload)
}

// Create subgroup
func (g *gitlab) CreateSubgroup(name string, path string, group *int) (int, error) {
	// Inherit some attributes from the parent group
	parentGroupDetail, err := g.viewGroup(*group)
	if err != nil {
		return 0, err
	}

	// Create the POST request payload
	payload := gitlabCreateSubgroupRequest{
		Name:                  name,
		Path:                  path,
		ParentID:              group,
		Visibility:            parentGroupDetail.Visibility,
		RequestAccessEnabled:  parentGroupDetail.RequestAccessEnabled,
		ProjectCreationLevel:  "maintainer",
		SubgroupCreationLevel: "owner",
	}

	return g.createGroup(payload)
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
					*channel <- fmt.Sprintf(messageErrorSetDefaultBranch, defaultBranch, projectID, err.Error())
				} else {
					*channel <- fmt.Sprintf(messageSetDefaultBranch, defaultBranch, projectID)
				}
			case 6, 3:
				err = g.setDefaultBranch(projectID, "staging")
				if err != nil {
					*channel <- fmt.Sprintf(messageErrorSetDefaultBranch, "staging", projectID, err.Error())
				} else {
					*channel <- fmt.Sprintf(messageSetDefaultBranch, "staging", projectID)
				}
			case 1:
				err = g.setDefaultBranch(projectID, "develop")
				if err != nil {
					*channel <- fmt.Sprintf(messageErrorSetDefaultBranch, "develop", projectID, err.Error())
				} else {
					*channel <- fmt.Sprintf(messageSetDefaultBranch, "develop", projectID)
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

func NewGitlab(apiURL string, token string, mirror GitlabMirrorOptions, exclusions GitlabExclusionsConfig) Gitlab {
	return &gitlab{
		apiURL:     apiURL,
		token:      token,
		mirror:     mirror,
		exclusions: exclusions,
	}
}
