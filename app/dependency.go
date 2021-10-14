package app

import (
	"strings"
)

type DependencyType string

const (
	Mutable   DependencyType = "mutable"
	Immutable DependencyType = "immutable"
)

type Dependency struct {
	Instances     []struct{
		ID string `json:"instanceid"`
		Version string `json:"version"`
	}
	Configuration interface{}
	Link          string
	Type          DependencyType
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
