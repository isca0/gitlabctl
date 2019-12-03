package model

// Group binds in json the received group response.
type Group struct {
	ID                             int    `json:"id"`
	Name                           string `json:"name"`
	WebURL                         string `json:"web_url"`
	Path                           string `json:"path"`
	Description                    string `json:"description"`
	Visibility                     string `json:"visibility"`
	ShareWithGroupLock             bool   `json:"share_with_group_lock"`
	RequireTwoFactorAuthentication bool   `json:"require_two_factor_authentication"`
	TwoFactorGracePeriod           int    `json:"two_factor_grace_period"`
	ProjectCreationLevel           string `json:"project_creation_level"`
	LfsEnabled                     bool   `json:"lfs_enabled"`
	AvatarURL                      string `json:"avatar_url"`
	RequestAccessEnabled           bool   `json:"request_access_enabled"`
	FullName                       string `json:"full_name"`
	FullPath                       string `json:"full_path"`
	ParentID                       int    `json:"parent_id"`
}

// GroupData is the group creation used on http.Header.
type GroupData struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`
	AvatarURL   string `json:"avatar_url"`
	ParentID    int    `json:"parent_id"`
}
