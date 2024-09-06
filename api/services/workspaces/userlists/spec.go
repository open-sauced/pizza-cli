package userlists

import (
	"time"

	"github.com/open-sauced/pizza-cli/api/services"
)

type DbUserListContributor struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	ListID    string    `json:"list_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type DbUserList struct {
	ID           string                  `json:"id"`
	UserID       int                     `json:"user_id"`
	Name         string                  `json:"name"`
	IsPublic     bool                    `json:"is_public"`
	IsFeatured   bool                    `json:"is_featured"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
	DeletedAt    *time.Time              `json:"deleted_at"`
	Contributors []DbUserListContributor `json:"contributors"`
}

type GetUserListsResponse struct {
	Data []DbUserList      `json:"data"`
	Meta services.MetaData `json:"meta"`
}

type CreatePatchUserListRequest struct {
	Name         string                             `json:"name"`
	IsPublic     bool                               `json:"is_public"`
	Contributors []CreateUserListRequestContributor `json:"contributors"`
}

type CreateUserListRequestContributor struct {
	Login string `json:"login"`
}

type CreateUserListResponse struct {
	ID          string `json:"id"`
	UserListID  string `json:"user_list_id"`
	WorkspaceID string `json:"workspace_id"`
}
