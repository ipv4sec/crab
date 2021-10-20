package status

type Status struct {
	PK        uint   `gorm:"primaryKey" json:"-"`
	ID        string `json:"id"`
	Component string `json:"component"`
	Status    int    `json:"status"` // 0 默认状态. 1 已部署
}

func (Status) TableName() string {
	return "t_status"
}
