package models

import "time"

type Status struct {
	ID        int        `json:"id"`
	ProjectID int        `json:"project_id"`
	Label     string     `json:"label"`
	Color     string     `json:"color"`
	Position  int        `json:"position"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
