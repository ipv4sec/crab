package cache

import (
	"crab/db"
	"k8s.io/klog/v2"
	"time"
)

type Node struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Data string
}

func (receiver *Node) Query(data string) string  {
	var err error
	var node Node
	err = db.Client.Model(Node{}).Find(&node).Where("data = ?", data).Error
	if err != nil {
		klog.Errorln(err)
	}
	return node.Data
}

func (receiver *Node) Add(data string)  {

}
