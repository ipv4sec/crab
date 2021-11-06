package utils

import (
	"archive/zip"
	"crab/aam/v1alpha1"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsTrait(s []struct {
	Type       string `yaml:"type"`
	Properties v1alpha1.Properties `yaml:"properties"`
}, e string) bool {
	for _, a := range s {
		if strings.Contains(a.Type, e) {
			return true
		}
	}
	return false
}

func MergeStringMaps(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func UnZip(dst, src string) error {
	zr, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zr.Close()
	if dst != "" {
		if err := os.MkdirAll(dst, 0755); err != nil {
			return err
		}
	}
	for _, file := range zr.File {
		if strings.HasPrefix(file.Name, "__MACOSX") {
			continue
		}
		path := filepath.Join(dst, file.Name)
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
			continue
		}
		fr, err := file.Open()
		if err != nil {
			return err
		}
		fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		_, err = io.Copy(fw, fr)
		if err != nil {
			return err
		}
		_ = fw.Close()
		_ = fr.Close()
	}
	return nil
}
