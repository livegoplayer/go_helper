package helper

import (
	"os"
	"path/filepath"
)

//用来存储文件目录相关帮助函数
//转换目录分隔符为对应系统的
func PathToCommon(str string) string {
	return filepath.FromSlash(str)
}

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
