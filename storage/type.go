package storage

import (
	"crab/sys"
	"time"
)

type DiskData struct {
	Value []*sys.LocalDisk
	UpdateAt time.Time
}

type AddrData struct {
	Value []string
	UpdateAt time.Time
}