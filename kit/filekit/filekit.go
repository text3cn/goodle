package filekit

import (
	"os"
	"path/filepath"
)

// 获取一个绝对路径所属目录
func Dir(absolutePath string) string {
	return filepath.Dir(absolutePath)
}

// 判断文件/文件夹是否存在
func PathExists(absolutePath string) (bool, error) {
	_, err := os.Stat(absolutePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 判断是否是文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断是否是文件而不是目录
func IsFile(path string) bool {
	return !IsDir(path)
}
