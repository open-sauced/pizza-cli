package workspaces

import (
	"time"

	"github.com/open-sauced/pizza-cli/api/services"
)

type DbWorkspace struct {
	ID          string              `json:"id"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	DeletedAt   *time.Time          `json:"deleted_at"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	IsPublic    bool                `json:"is_public"`
	PayeeUserID *int                `json:"payee_user_id"`
	Members     []DbWorkspaceMember `json:"members"`
}

type DbWorkspaceMember struct {
	ID          string     `json:"id"`
	UserID      int        `json:"user_id"`
	WorkspaceID string     `json:"workspace_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	Role        string     `json:"role"`
}

type DbWorkspacesResponse struct {
	Data []DbWorkspace     `json:"data"`
	Meta services.MetaData `json:"meta"`
}

type CreateWorkspaceRequestRepoInfo struct {
	FullName string `json:"full_name"`
}

type CreateWorkspaceRequest struct {
	Name         string                           `json:"name"`
	Description  string                           `json:"description"`
	Members      []string                         `json:"members"`
	Repos        []CreateWorkspaceRequestRepoInfo `json:"repos"`
	Contributors []string                         `json:"contributors"`
}
