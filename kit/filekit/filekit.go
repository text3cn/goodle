package filekit

import (
	"fmt"
	"github.com/text3cn/goodle/kit/strkit"
	"io"
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

// 根据文件名获取后缀
func GetSuffix(filename string) string {
	s := strkit.Explode(".", filename)
	last := s[len(s)-1]
	if last == filename {
		return ""
	}
	return last
}

// 创建目录
func MkDir(dir string, mode os.FileMode) {
	var err error
	if _, err = os.Stat(dir); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(dir, mode); err != nil {
			fmt.Println("创建目录失败：" + dir)
		}
	}
}

// 扫描指定目录下所有文件及文件夹
// dir 指定要扫描的目录
// return：
// files 文件数组
// dirs  文件夹数组
func Scandir(dir string) (files []string, dirs []string) {
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			// 扫描到的 ./ (自身) 不要
			if path != dir {
				dirs = append(dirs, path)
			}
			return nil
		}
		files = append(files, path)
		return nil
	})
	return
}

// 创建文件，覆盖创建
func createFile(filepath string, content string) {
	f, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		_, err = f.Write([]byte(content))
	}
}

// 读取文件内容
func readFile(filepath string) string {
	ret := ""
	f, err := os.OpenFile(filepath, os.O_RDONLY, 0600)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		contentByte, _ := io.ReadAll(f)
		ret = string(contentByte)
	}
	return ret
}
