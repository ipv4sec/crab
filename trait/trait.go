package trait

import "time"

type Trait struct {
	PK    int64 `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Ver   string `json:"apiVersion"`
	Value string `json:"value"`

	Type int
	// 0 内置, 系统级别, 不可删
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
