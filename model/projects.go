package model

import "time"

// Projects satisfies every project methods.
type Projects struct {
	ID                int           `json:"id"`
	Description       string        `json:"description"`
	Name              string        `json:"name"`
	Visibility        string        `json:"visibility"`
	NameWithNamespace string        `json:"name_with_namespace"`
	Path              string        `json:"path"`
	PathWithNamespace string        `json:"path_with_namespace"`
	CreatedAt         time.Time     `json:"created_at"`
	DefaultBranch     string        `json:"default_branch"`
	TagList           []interface{} `json:"tag_list"`
	SSHURLToRepo      string        `json:"ssh_url_to_repo"`
	HTTPURLToRepo     string        `json:"http_url_to_repo"`
	WebURL            string        `json:"web_url"`
	AvatarURL         string        `json:"avatar_url"`
	StarCount         int           `json:"star_count"`
	ForksCount        int           `json:"forks_count"`
	LastActivityAt    time.Time     `json:"last_activity_at"`
	Namespace         Namespace     `json:"namespace"`
	Custom            CustomFlags   `json:"custom_flags"`
}

// Namespace binds the json received from endpoint project.
type Namespace struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Kind     string `json:"kind"`
	FullPath string `json:"full_path"`
	ParentID int    `json:"parent_id"`
}

// CustomFlags brigs custom paramethers to projects.
type CustomFlags struct {
	ClonePath string `json:"clone_path"`
	BareRepo  bool   `json:"bare_repo"`
	NewRepo   string `json:"new_repo"`
}
