package app

import "time"

type App struct {
	PK     int64  `gorm:"primaryKey" json:"-"`
	ID     string `json:"id"`
	Status int    `json:"status"`

	Name           string `json:"name"`
	Version        string `json:"version"`
	Configurations string `json:"configurations"`
	Dependencies   string `json:"dependencies"`

	Manifest string `json:"manifest"`

	Parameters string `json:"parameters"`
	Deployment string `gorm:"column:deployment" json:"deployment"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (App) TableName() string {
	return "t_app"
}
