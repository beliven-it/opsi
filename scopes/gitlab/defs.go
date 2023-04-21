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
	CreateProject(string, string, int, string, bool) (int, error)
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
