package app

import "time"

type App struct {
	ID            string `gorm:"primaryKey" json:"-"`
	Name          string `json:"name"`
	Version       string `json:"version"`
	Status        int    `json:"status"` // 0 正在部署中 1 部署完成 2 卸载中 3 卸载完成
	Namespace     string `json:"namespace"`
	Manifest      string `json:"manifest"`
	Dependencies  string `json:"dependencies"`
	Parameters    string `json:"parameters"`
	Configuration string

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Deployment string `gorm:"column:deployment"`
}

func (App) TableName() string {
	return "app"
}

type Manifest struct {
	Name                 string                 `json:"name"`
	Version              string                 `json:"version"`
	Services             map[string]interface{} `json:"services"`
	Dependencies         interface{}
	RuntimeConfiguration interface{} `json:"userconfig"`
}
