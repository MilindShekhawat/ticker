package models

const (
	RoleOwner  = 1
	RoleAdmin  = 2
	RoleMember = 3
	RoleViewer = 4
)

type ProjectMember struct {
	ProjectID int `json:"project_id"`
	UserID    int `json:"user_id"`
	Role      int `json:"role"`
}
