package utils

const (
	// 01 表示代码本身的错误
	ErrBadRequest      = 10101
	ErrInternalServer  = 10102

	// 02 表示数据库的错误
	ErrDatabaseBadRequest = 10201
	ErrDatabaseInternalServer = 10201

	// 03 表示集群本身的错误
	ErrClusterBadRequest = 10301
	ErrClusterInternalServer = 10301
)