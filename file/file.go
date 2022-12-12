package fileutil

import (
	"os"
	"path/filepath"
	"strings"
)

// GetDirFiles 获取目录下的文件夹
func GetDirFiles(dir string) []string {
	//将地址统一转化为斜杠
	//ToSlash将所有反斜杠转成斜杠；FromSlash将所有斜杠转成反斜杠
	dir = filepath.ToSlash(dir)

	var files []string                                     //存储文件夹名
	readDir, err := os.ReadDir(dir)                        //获取目录下所有文件
	if err != nil || readDir == nil || len(readDir) <= 0 { //获取错误 || 获取为空 || 文件数为0
		return files
	}

	//遍历所有文件
	for _, entry := range readDir {
		if entry.IsDir() { //判断是否为文件夹
			files = append(files,
				strings.TrimRight(dir, "/")+ //如果最右边有”/“，则去掉，否则忽略
					"/"+entry.Name())
		}
	}
	return files
}
