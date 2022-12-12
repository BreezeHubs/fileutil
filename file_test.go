package main

import (
	fileutil "fileutil/file"
	"testing"
)

func Test_GetDirFiles(t *testing.T) {
	t.Log(fileutil.GetDirFiles("../"))
}

func Test_FastCopyFile(t *testing.T) {
	file, err := fileutil.FastCopyFile("./README.md", "./README.md.bak")
	t.Log(file, err)
}
