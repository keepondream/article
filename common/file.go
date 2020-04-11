package common

import (
	"fmt"
	"os"
	"syscall"
)

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// 创建目录
func Mkdir(path string) bool {
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)
	err := os.MkdirAll(path, 0777)
	fmt.Println(err)
	if err != nil {
		panic(err)
	}
	fmt.Println("创建成功")
	return true
}

// 获取项目根目录
func GetBasePath() string {
	//获取当前工作目录的根路径
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path
}
