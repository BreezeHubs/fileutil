package fileutil

import (
	"io"
	"os"
)

// CopyFile 拷贝文件
func CopyFile(src, dest string) (int64, error) {
	srcFile, err := os.Open(src) //打开源文件
	if err != nil {
		return 0, err
	}
	defer srcFile.Close() //用完关闭

	//os.O_CREATE可创建；os.O_APPEND可追加；os.O_WRONLY可读写
	//当前用户：6可读可写；其他用户：4只可读
	destFile, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644) //打开或创建拷贝文件，可创建可读写
	if err != nil {
		return 0, err
	}
	defer destFile.Close() //用完关闭

	return io.Copy(destFile, srcFile) //使用io库进行复制文件
}

// FastCopyFile 拷贝文件
func FastCopyFile(src, dest string) (int64, error) {
	srcFile, err := os.Open(src) //打开源文件
	if err != nil {
		return 0, err
	}
	defer srcFile.Close() //用完关闭

	destFile, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644) //打开或创建拷贝文件，可创建可读写
	if err != nil {
		return 0, err
	}
	defer destFile.Close() //用完关闭

	//利用缓冲池 边读边写
	buf := make([]byte, 1024)
	var count int //统计字节数
	for {
		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}
		if n == 0 {
			break
		}
		destFile.Write(buf[:n])
		count += n //累加字节数
	}
	return int64(count), nil
}
