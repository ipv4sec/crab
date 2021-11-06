package status

import "time"

type Status struct {
	PK      uint   `gorm:"primaryKey" json:"-"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Status  int    `json:"status"` // 0 默认状态. 1 已部署
	Message string `json:"message"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Status) TableName() string {
	return "t_status"
}
