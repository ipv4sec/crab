package app

import (
	"strings"
)

type DependencyType string

const (
	Mutable   DependencyType = "mutable"
	Immutable DependencyType = "immutable"
)

type Instance struct {
	ID      string `json:"id"`
	Name string `json:"name"`
}

type Dependency struct {
	Instances      []Instance     `json:"instances"`
	Configurations struct{}       `json:"userconfigs"`
	Link           string         `json:"location"`
	Type           DependencyType `json:"type"`
}

func Link(value string) (DependencyType, string) {
	if "user-defined" == value {
		return Mutable, ""
	}
	if strings.HasPrefix(value, "user-defined(") &&
		strings.HasSuffix(value, ")") {
		return Mutable, strings.TrimSuffix(strings.TrimPrefix(value, "user-defined("), ")")
	}
	return Immutable, value
}
