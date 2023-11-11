package main

import (
	"os"
	"path/filepath"
)

func Exists(pathLike string) bool {
	_, err := os.Stat(pathLike); 
	return !os.IsNotExist(err)
}

func IsFile(pathLike string) bool {
	if !Exists(pathLike) {
		return false
	}

	info, _ := os.Stat(pathLike)
	return !info.IsDir()
}

func IsDir(pathLike string) bool {
	if !Exists(pathLike) {
		return false
	}
	return !IsFile(pathLike)
}

func IsFileReadable(filePath string) bool {
	if !IsFile(filePath) {
		return false
	}

	_, err := os.Open(filePath)
	return err == nil
}

func IsDirWritable(dirPath string) bool {
	if !IsFile(dirPath) {
		return false
	}

	info, _ := os.Stat(dirPath)
	return info.Mode().Perm() & 0200 != 0
}

func AbsolutePath(pathLike string) string {
	if !Exists(pathLike) {
		panic("Path does not exist!!!")
	}
	absPath, _ := filepath.Abs(pathLike)
	return absPath
}