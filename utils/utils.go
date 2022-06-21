package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func FindFile(path string, patterns []string) (string, error) {
	foundPath := ""
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			for _, pattern := range patterns {
				if strings.Contains(filepath.Base(path), pattern) {
					foundPath = path
					return io.EOF
				}
			}
		}
		return nil
	})
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return foundPath, err
	}
	return foundPath, nil
}
func findFiles(path string, patterns []string) (string, error) {
	foundPath := ""
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			for _, pattern := range patterns {
				if strings.Contains(filepath.Base(path), pattern) {
					foundPath = path
					return io.EOF
				}
			}
		}
		return nil
	})
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return foundPath, err
	}
	return foundPath, nil
}

func Copy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	dstFileStat, _ := os.Stat(dst)
	if dstFileStat.IsDir() {
		dst = filepath.Join(dst, filepath.Base(src))
	}
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	_, err = io.Copy(destination, source)
	return err
}
