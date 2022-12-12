# 关于文件操作的包

## 1 file包
- 获取目录下的文件夹
    ```go
    dirs := fileutil.GetDirFiles("../")
    fmt.Println(dirs)
    ```
- 拷贝文件
    ```go
    file, err := fileutil.CopyFile("./README.md", "./README.md.bak") //全量拷贝
	fmt.Println(file, err)
    file, err = fileutil.FastCopyFile("./README.md", "./README.md.bak") //缓冲池边读边拷贝
	fmt.Println(file, err)
    ```

## 2 log包
- 暂未更新

## 3 yaml包
- 解析yaml文件
    ```go
    config, _ := yaml.LoadConfig("./config.yaml") //读取yaml配置文件
    getString := config.GetString("devices[0].nodes[0].index") //获取key为string格式
    getInt := config.GetInt("ipport") //获取key为int格式
    fmt.Println(getString, getInt)
    ```