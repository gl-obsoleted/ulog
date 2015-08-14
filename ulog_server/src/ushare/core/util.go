package core

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func DirExists(dir string) bool {
	// check if the source dir exist
	src, err := os.Stat(dir)
	if err != nil {
		return false
	}

	if !src.IsDir() {
		return false
	}

	return true
}

func FileExists(file string) bool {
	src, err := os.Stat(file)
	if err != nil {
		return false
	}

	if src.Size() == 0 {
		return false
	}

	return true
}

func GetFileSize(file string) int64 {
	src, err := os.Stat(file)
	if err != nil {
		return 0
	}

	return src.Size()
}

func GetFileMD5(file string) string {
	_, err := os.Stat(file)
	if err != nil {
		return ""
	}

	f, inerr := os.Open(file)
	if inerr != nil {
		return ""
	}

	md5h := md5.New()
	io.Copy(md5h, f)
	return fmt.Sprintf("%x", md5h.Sum([]byte(""))) //md5
}

func CreateDirIfNotExists(dir string) error {
	if !DirExists(dir) {
		if err := os.MkdirAll(dir, os.ModeDir); err != nil {
			return err
		}
	}

	return nil
}
