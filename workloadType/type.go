package workloadType

import "time"

type WorkloadType struct {
	Id     int64  `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Ver   string `json:"apiVersion" gorm:"column:ver"`
	Value string `json:"value"`

	Type int `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (WorkloadType) TableName() string {
	return "t_type"
}
