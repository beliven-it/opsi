package gitlab

const gitlabOwnerPermission int = 50
const gitlabMaintainerPermission int = 40
const gitlabDeveloperPermission int = 30
const gitlabDefaultGroupMemberMaintainer string = "default_group_member_maintainer"
const gitlabDefaultGroupMemberDeveloper string = "default_group_member_developer"
const gitlabDefaultGroupMemberOwner string = "default_group_member_owner"
const gitlabDefaultGroupMember string = "default_group_member"

type gitlab struct {
	token  string
	apiURL string
	mirror GitlabMirrorOptions
}

type Gitlab interface {
	CreateEnvs(string, string, string) error
	ListEnvs(string, string) error
	DeleteEnvs(string, string) error
	CreateProject(ProjectRequest) (int, error)
	CreateSubgroup(string, string, *int) (int, error)
	CreateGroup(string, string, string) (int, error)
	BulkSettings(*chan string) error
	Deprovionioning(string) error
}

type GitlabMirrorOptions struct {
	Token     string `mapstructure:"token"`
	ApiURL    string `mapstructure:"api_url"`
	GroupPath string `mapstructure:"group_path"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	GroupID   int    `mapstructure:"group_id"`
}

type gitlabCreateMirrorRequest struct {
	Enabled               bool   `json:"enabled"`
	URL                   string `json:"url"`
	OnlyProtectedBranched bool   `json:"only_protected_branches"`
}

type gitlabDefaultUser struct {
	tipology   string
	permission int
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
	ID                   int    `json:"id"`
	Visibility           string `json:"visibility"`
	RequestAccessEnabled bool   `json:"request_access_enabled"`
}

type gitlabCreateProjectRequest struct {
	Name                             string `json:"name"`
	Path                             string `json:"path"`
	Visibility                       string `json:"visibility"`
	NamespaceID                      int    `json:"namespace_id"`
	MergeMethod                      string `json:"merge_method"`
	LFSEnabled                       bool   `json:"lfs_enabled"`
	SharedRunnersEnabled             bool   `json:"shared_runners_enabled"`
	InitializeWithReadME             bool   `json:"initialize_with_readme"`
	SquashOption                     string `json:"squash_option"`
	PackagesEnabled                  bool   `json:"packages_enabled"`
	MirrorTriggerBuilds              bool   `json:"mirror_trigger_builds"`
	BuildsAccessLevel                string `json:"builds_access_level"`
	AnalyticsAccessLevel             string `json:"analytics_access_level"`
	PagesAccessLevel                 string `json:"pages_access_level"`
	ContainerRegistryAccessLevel     string `json:"container_registry_access_level"`
	OperationsAccessLevel            string `json:"operations_access_level"`
	IssuesAccessLevel                string `json:"issues_access_level"`
	MergeRequestAccessLevel          string `json:"merge_request_access_level"`
	ReleasesAccessLevel              string `json:"releases_access_level"`
	EnvironmentsAccessLevel          string `json:"environments_access_level"`
	FeatureFlagsAccessLevel          string `json:"feature_flags_access_level"`
	MonitorAccessLevel               string `json:"monitor_access_level"`
	RepositoryAccessLevel            string `json:"repository_access_level"`
	RequirementsAccessLevel          string `json:"requirements_access_level"`
	InfrastrucureAccessLevel         string `json:"infrastructure_access_level"`
	SecurityAndComplianceAccessLevel string `json:"security_and_compliance_access_level"`
	SnippetAccessLevel               string `json:"snippet_access_level"`
	WikiAccessLevel                  string `json:"wiki_access_level"`
	ForkingAccessLevel               string `json:"forking_access_level"`
}

type gitlabAddUserToGroupRequest struct {
	ID          int `json:"id"`
	UserID      int `json:"user_id"`
	AccessLevel int `json:"access_level"`
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
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Default bool   `json:"default"`
}

type gitlabSetupBranchRequest struct {
	Name             string `json:"name"`
	PushAccessLevel  int    `json:"push_access_level"`
	MergeAccessLevel int    `json:"merge_access_level"`
}

type gitlabUser struct {
	ID   int    `json:"id"`
	Note string `json:"note"`
}

type gitlabCreateEnvRequest struct {
	VariableType     string `json:"variable_type"`
	Key              string `json:"key"`
	Value            string `json:"value"`
	EnvironmentScope string `json:"environment_scope"`
	Masked           bool   `json:"masked"`
	Protected        bool   `json:"protected"`
}

type ProjectRequest struct {
	Name          string
	Path          string
	Visibility    string
	DefaultBranch string
	Mirror        bool
	SharedRunners bool
	Group         int
}

var defaultGitlabCreatePayload = gitlabCreateProjectRequest{
	MergeMethod:                      "ff",
	Visibility:                       "private",
	LFSEnabled:                       false,
	SharedRunnersEnabled:             false,
	InitializeWithReadME:             true,
	SquashOption:                     "never",
	PackagesEnabled:                  true,
	IssuesAccessLevel:                "disabled",
	PagesAccessLevel:                 "disabled",
	OperationsAccessLevel:            "disabled",
	AnalyticsAccessLevel:             "disabled",
	SecurityAndComplianceAccessLevel: "disabled",
	SnippetAccessLevel:               "enabled",
	WikiAccessLevel:                  "disabled",
	ForkingAccessLevel:               "disabled",
}

var defaultGitlabMirrorCreatePayload = gitlabCreateProjectRequest{
	MergeMethod:                      "ff",
	Visibility:                       "private",
	LFSEnabled:                       false,
	SharedRunnersEnabled:             false,
	InitializeWithReadME:             true,
	SquashOption:                     "never",
	PackagesEnabled:                  false,
	BuildsAccessLevel:                "disabled",
	AnalyticsAccessLevel:             "disabled",
	PagesAccessLevel:                 "disabled",
	ContainerRegistryAccessLevel:     "disabled",
	OperationsAccessLevel:            "disabled",
	IssuesAccessLevel:                "disabled",
	MergeRequestAccessLevel:          "disabled",
	ReleasesAccessLevel:              "disabled",
	EnvironmentsAccessLevel:          "disabled",
	FeatureFlagsAccessLevel:          "disabled",
	MonitorAccessLevel:               "disabled",
	RepositoryAccessLevel:            "enabled",
	RequirementsAccessLevel:          "disabled",
	InfrastrucureAccessLevel:         "disabled",
	SecurityAndComplianceAccessLevel: "disabled",
	SnippetAccessLevel:               "disabled",
	WikiAccessLevel:                  "disabled",
	ForkingAccessLevel:               "disabled",
}

const projectEndpoint = "/projects"

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

var defaultProtectedTags = map[string]interface{}{
	"allowed_to_create": []map[string]interface{}{
		{
			"access_level": gitlabMaintainerPermission,
		},
	},
	"name": "*",
}

const messageSetDefaultBranch = "Set default branch %s for project #%d"
const messageErrorSetDefaultBranch = "Error on set default branch %s for project #%d: %s"
