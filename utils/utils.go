package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
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

func GetOwner(file string) (int, int, error) {
	info, _ := os.Stat(file)
	uid := 0
	gid := 0
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		uid = int(stat.Uid)
		gid = int(stat.Gid)
	} else {
		// we are not in linux, this won't work anyway in windows,
		// but maybe you want to log warnings
		fmt.Println("fail to get uid")
	}
	return uid, gid, nil
}

func Copy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("fail to stat %s: %s", src, err)
	}
	dstFileStat, err := os.Stat(dst)
	if err != nil {
		return fmt.Errorf("fail to stat %s: %s", dst, err)
	}
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
	stat, err := os.Stat(src)
	uid, gid, err := GetOwner(src)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = os.Chmod(dst, stat.Mode())
	err = os.Chown(dst, uid, gid)
	return err
}
