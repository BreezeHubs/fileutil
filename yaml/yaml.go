package yaml

import (
	"os"
	"strings"
	"unsafe"
)

// Yaml 预解析
type Yaml struct {
	data map[string]any
}

// 原始解析数据结构
type item struct {
	level  int    //层级
	val    string //值
	key    string //键
	cIndex []int  //子数据索引
}

// LoadConfig 加载配置文件，解析
func LoadConfig(file string) (*Yaml, error) {
	//读取目标文件
	readFile, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	content := *(*string)(unsafe.Pointer(&readFile)) //解析成文本
	list := strings.Split(content, "\n")             //根据换行符分割成数组

	listItems := make([]*item, 0) //存放解析数据
	for i := 0; i < len(list); i++ {
		level, key, val := getListItem(list[i]) //分割层级、key、value
		if key == "" {                          //如果当前为空行，则忽略
			continue
		}

		//判断为数组
		if strings.HasPrefix(key, "-") {
			//如果该数组不是第一层
			if level != 0 {
				//需要找它的上层所在行数，-1为找不到上层
				if pIndex := getParentLineIndex(listItems, level); pIndex != -1 {
					//能找到上层行数时，把当前数据行数（索引）存到上层的child层数据中
					listItems[pIndex].cIndex = append(listItems[pIndex].cIndex, len(listItems))
				}
			}

			//存入解析数据
			listItems = append(listItems, &item{level: level}) //数组只存层级信息
			level += 2                                         //数组本身占一层，所以+2到数组的层
			key = strings.Replace(key, "- ", "", 1)            //数组数组的key
		}

		//非数组
		//如果该数组不是第一层
		if level != 0 {
			//需要找它的上层所在行数，-1为找不到上层
			if pIndex := getParentLineIndex(listItems, level); pIndex != -1 {
				//能找到上层行数时，把当前数据行数（索引）存到上层的child层数据中
				listItems[pIndex].cIndex = append(listItems[pIndex].cIndex, len(listItems))
			}
		}
		//存入解析数据 非数组
		listItems = append(listItems, &item{level: level, val: val, key: key})
	}

	//debug list
	//for _, listItem := range listItems {
	//	fmt.Printf("%+v\n\n", listItem)
	//}

	//将解析的数据分层存储
	m := getResult(listItems)
	//fmt.Printf("%v", m)
	return &Yaml{data: m}, nil
}

// 寻找对应层级的上层所在行数
// list当前所有层级数据、level该数据层级
// 返回上层所在行数，上层寻找失败返回-1
func getParentLineIndex(list []*item, level int) int {
	for i := len(list) - 1; i >= 0; i-- {
		if list[i].level < level { //当出现层级小于该数据层级时，确定为该数据的上层
			return i
		}
	}
	return -1
}

// 分割层级、key、value
func getListItem(s string) (level int, key, val string) {
	level = len(s) - len(strings.TrimLeft(s, " ")) //计算空格
	s = strings.Trim(s, " ")                       //去掉空格
	idx := strings.Index(s, ":")                   //查找key和value分隔符位置

	//可能存在空行，判断不为空
	if idx != -1 {
		key = s[:idx]                                    //得到key
		val = valueTrim(s[idx+1:], " ", "\r", "\"", "'") //value去掉空格、换行、双引号、单引号
	}
	return
}

// trim过滤字符串
func valueTrim(s string, trim ...string) string {
	if len(trim) > 0 {
		for _, t := range trim {
			s = strings.Trim(s, t)
		}
	}
	return s
}

// 将解析的数据分层存储
func getResult(list []*item) map[string]any {
	m := make(map[string]any)

	//逐级递归
	for i, item := range list {
		if item.level == 0 {
			m[item.key] = getChild(i, list)
		}
	}
	return m
}

// 递归下级
func getChild(idx int, list []*item) any {
	it := list[idx] //获取当前数据
	val := it.val   //取到当前value

	//如果没有下级，表明单级数据，直接返回该值
	if len(it.cIndex) == 0 {
		return val
	}

	//存在下级
	m := make(map[string]any) //如果下级是map
	slice := make([]any, 0)   //如果下级是数组
	for _, index := range it.cIndex {
		it := list[index]            //获取下级数据
		res := getChild(index, list) //递归 获取下级的数据

		if it.key == "" {
			//下级是数组
			slice = append(slice, res)
		} else {
			//下级是map
			m[it.key] = res
		}
	}

	if len(m) > 0 { //map
		return m
	} else { //数组
		return slice
	}
}
