package yaml

import (
	"strconv"
	"strings"
)

// Get 获取key索引出的值
func (y *Yaml) Get(key string) any {
	var (
		keys = getKeys(key)          //解析key索引
		ok   bool                    //判断是否正常类型转换
		tmp  any            = y.data //将数据预存下来
	)
	for _, k := range keys {
		//判断是否能转成 map[string]any
		if _, ok = tmp.(map[string]any); !ok {
			return ""
		}

		//判断map是否能取到对应下标
		tmp, ok = tmp.(map[string]any)[k.key]
		if !ok {
			return ""
		}

		//是否存在数组下标
		if k.hasIndex {
			//判断是否能转成 []any
			if _, ok = tmp.([]any); !ok {
				return ""
			}

			//判断是否超过最大索引
			if len(tmp.([]any))-1 < k.index {
				return ""
			}

			tmp = tmp.([]any)[k.index] //取出索引值
		}
	}
	//判断是否能取到值
	if _, ok = tmp.(string); !ok {
		return ""
	}
	return tmp
}

func (y *Yaml) GetInt(key string) int {
	s := y.GetString(key)
	i, _ := strconv.Atoi(s)
	return i
}

func (y *Yaml) GetString(key string) string {
	return y.Get(key).(string)
}

// 解析的key数据结构
type keyItem struct {
	key      string //key
	index    int    //数组下标
	hasIndex bool   //是否存在数组下标
}

// 解析如 devices[0].nodes[0].index 的key
func getKeys(key string) []*keyItem {
	keys := make([]*keyItem, 0)

	list := strings.Split(key, ".") //分割 devices[0].nodes[0].index => [devices[0], nodes[0], index]
	for _, k := range list {
		listK := strings.Split(k, "[")   //devices[0] => [devices, 0]]
		kItem := &keyItem{key: listK[0]} //存入key

		//如果存在数组下标 如：devices[0] => [devices, 0]]
		if len(listK) > 1 {
			kItem.hasIndex = true
			kItem.index, _ = strconv.Atoi( //数组下标转成int，可以忽略错误
				strings.TrimRight(listK[1], "]"), //去掉未处理的 ] 符号
			)
		}
		keys = append(keys, kItem)
	}
	return keys
}
