package app

import "time"

type App struct {
	ID           int64 `gorm:"primaryKey" json:"-"`
	UUID         string `json:"namespace"`
	Status       int    `json:"status"` // TODO 0 正在部署中 1 部署完成 2 卸载中 3 卸载完成

	Name         string `json:"name"`
	Version      string `json:"version"`
	Manifest     string `json:"manifest"`

	//Dependencies string `json:"dependencies"`
	//Parameters   string `json:"parameters"`
	//Config       string
	Deployment     string `gorm:"column:deployment"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (App) TableName() string {
	return "t_app"
}