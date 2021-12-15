package workloadType

import "time"

type WorkloadType struct {
	PK     int64  `gorm:"primaryKey" json:"-"`
	Name  string `json:"name"`
	Ver   string `json:"apiVersion"`
	Value string `json:"value"`

	Type int
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
