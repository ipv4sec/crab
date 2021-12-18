package app

import (
	"time"
)

type App struct {
	PK     int64  `gorm:"primaryKey" json:"-"`
	ID     string `json:"id"`

	Name           string `json:"name"`
	Version        string `json:"version"`
	Configurations string `json:"configurations"`
	Dependencies   string `json:"dependencies"`

	Manifest string `json:"manifest"`
	Entry string `json:"entry"`

	Parameters string `json:"parameters"`
	Additional string `json:"additional"`

	Deployment string `gorm:"column:deployment" json:"deployment"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (App) TableName() string {
	return "t_app"
}
