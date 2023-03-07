package onepassword

type onePassword struct {
	address string
	account OnePasswordAccount
}

type OnePassword interface {
	Deprovisioning(string) error
	Create(string) error
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
	ContentVersion int    `json:"content_version"`
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
