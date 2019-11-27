package model

// Groups binds in json the received groups response.
type Groups []struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	Path                 string `json:"path"`
	Description          string `json:"description"`
	Visibility           string `json:"visibility"`
	LfsEnabled           bool   `json:"lfs_enabled"`
	AvatarURL            string `json:"avatar_url"`
	WebURL               string `json:"web_url"`
	RequestAccessEnabled bool   `json:"request_access_enabled"`
	FullName             string `json:"full_name"`
	FullPath             string `json:"full_path"`
	ParentID             int    `json:"parent_id"`
}

// GroupCreation simplifies the groups creation.
type GroupCreation struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Description string `json:"description"`
	ParentID    int    `json:"parent_id"`
	Visibility  string `json:"visibility"`
}
