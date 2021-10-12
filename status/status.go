package status

type Status struct {
	ID        uint   `gorm:"primaryKey" json:"-"`
	Name      string `json:"name"`
	Component string `json:"component"`
	Status    int    `json:"status"`
}

func (Status) TableName() string {
	return "t_status"
}
