package utils

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

// PathToCommon 用来存储文件目录相关帮助函数
// 转换目录分隔符为对应系统的
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

func GetFileExtName(str string) string {
	fileSuffix := path.Ext(str)
	if len(fileSuffix) > 0 {
		return fileSuffix[1:]
	} else {
		return ""
	}
}

func GetCurPath() string {
	file, _ := exec.LookPath(os.Args[0])
	p, _ := filepath.Abs(file)
	rst := filepath.Dir(p)
	return rst
}

func GetFileExt(path string) string {
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
	}
	return ""
}
