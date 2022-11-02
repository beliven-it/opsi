package scopes

type postmarkEditRequest struct {
	SmtpApiActivated           bool   `json:"SmtpApiActivated"`
	RawEmailEnabled            bool   `json:"RawEmailEnabled"`
	BounceHookUrl              string `json:"BounceHookUrl"`
	PostFirstOpenOnly          bool   `json:"PostFirstOpenOnly"`
	TrackOpens                 bool   `json:"TrackOpens"`
	TrackLinks                 string `json:"TrackLinks"`
	IncludeBounceContentInHook bool   `json:"IncludeBounceContentInHook"`
	EnableSmtpApiErrorHooks    bool   `json:"EnableSmtpApiErrorHooks"`
}

type postmarkCreateRequest struct {
	postmarkEditRequest
	Name  string `json:"Name"`
	Color string `json:"Color"`
}

type postmarkServer struct {
	ID   int    `json:"ID"`
	Name string `json:"Name"`
}

type postmarkServersResponse struct {
	Servers []postmarkServer
}

type gitlabSubgroupRequest struct {
	Name                  string `json:"name"`
	Path                  string `json:"path"`
	ParentID              int    `json:"parent_id"`
	Visibility            string `json:"visibility"`
	ProjectCreationLevel  string `json:"project_creation_level"`
	SubgroupCreationLevel string `json:"subgroup_creation_level"`
	RequestAccessEnabled  bool   `json:"request_access_enabled"`
}

type gitlabSubgroupResponse struct {
	ID int `json:"id"`
}

type gitlabProjectRequest struct {
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
