package gitlab

type gitlabProject struct {
}

type gitlab struct {
	token          string
	groupID        int
	apiURL         string
	projectService gitlabProject
}

type Gitlab interface {
	CreateEnvs(string, string, string) error
	ListEnvs(string, string) error
	DeleteEnvs(string, string) error
	CreateProject(string, string, int) error
	CreateSubgroup(string, string, *int) error
	BulkSettings() error
	Deprovionioning(string) error
}

type gitlabCreateSubgroupRequest struct {
	Name                  string `json:"name"`
	Path                  string `json:"path"`
	ParentID              *int   `json:"parent_id"`
	Visibility            string `json:"visibility"`
	ProjectCreationLevel  string `json:"project_creation_level"`
	SubgroupCreationLevel string `json:"subgroup_creation_level"`
	RequestAccessEnabled  bool   `json:"request_access_enabled"`
}

type gitlabSubgroupResponse struct {
	ID int `json:"id"`
}

type gitlabCreateProjectRequest struct {
	Name                         string `json:"name"`
	Path                         string `json:"path"`
	NamespaceID                  int    `json:"namespace_id"`
	MergeMethod                  string `json:"merge_method"`
	AnalyticsAccessLevel         string `json:"analytics_access_level"`
	SecurityAndComplianceEnabled bool   `json:"security_and_compliance_enabled"`
	IssuesEnabled                bool   `json:"issues_enabled"`
	ForkingAccessLevel           string `json:"forking_access_level"`
	LFSEnabled                   bool   `json:"lfs_enabled"`
	WikiEnabled                  bool   `json:"wiki_enabled"`
	PagesAccessLevel             string `json:"pages_access_level"`
	OperationsAccessLevel        string `json:"operations_access_level"`
	SharedRunnersEnabled         bool   `json:"shared_runners_enabled"`
	InitializeWithReadME         bool   `json:"initialize_with_readme"`
	SquashOption                 string `json:"squash_option"`
}

type gitlabProjectResponse struct {
	ID int `json:"id"`
}

type gitlabProjectListVariable struct {
	VariableType     string `json:"variable_type"`
	Key              string `json:"key"`
	Value            string `json:"value"`
	Protected        bool   `json:"protected"`
	Masked           bool   `json:"masked"`
	Raw              bool   `json:"raw"`
	EnvironmentScope string `json:"environment_scope"`
}

type gitlabEntityWithID struct {
	ID int `json:"id"`
}

type gitlabBranchResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type gitlabSetupBranchRequest struct {
	Name             string `json:"name"`
	PushAccessLevel  int    `json:"push_access_level"`
	MergeAccessLevel int    `json:"merge_access_level"`
}

type gitlabUser struct {
	ID int `json:"id"`
}

type gitlabCreateEnvRequest struct {
	VariableType     string `json:"variable_type"`
	Key              string `json:"key"`
	Value            string `json:"value"`
	EnvironmentScope string `json:"environment_scope"`
	Masked           bool   `json:"masked"`
	Protected        bool   `json:"protected"`
}

var defaultGitlabCreatePayload = gitlabCreateProjectRequest{
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

const defaultBranch = "master"

var defaultProjectDevelopSettings = gitlabSetupBranchRequest{
	Name:             "develop",
	PushAccessLevel:  30,
	MergeAccessLevel: 30,
}

var defaultProjectStagingSettings = gitlabSetupBranchRequest{
	Name:             "staging",
	PushAccessLevel:  0,
	MergeAccessLevel: 30,
}

var defaultProjectMasterSettings = gitlabSetupBranchRequest{
	Name:             "master",
	PushAccessLevel:  0,
	MergeAccessLevel: 30,
}

var defaultCleanUpPolicy = map[string]interface{}{
	"container_expiration_policy_attributes": map[string]interface{}{
		"cadence":         "1month",
		"enabled":         true,
		"keep_n":          1,
		"older_than":      "14d",
		"name_regex":      ".*",
		"name_regex_keep": ".*-main",
	},
}
