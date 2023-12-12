package typekit

import (
	"encoding/json"
	"fmt"
)

// 序列化成json
func JsonEncode(val interface{}) []byte {
	json, err := json.Marshal(val)
	if err != nil {
		fmt.Println("Json encode error")
	}
	return json
}

// JSON反序列化
func JsonDecodeMap(str string) (thisMap map[string]interface{}, err error) {
	if err = json.Unmarshal([]byte(str), &thisMap); err != nil {
		fmt.Printf("Json反序列化为Map出错: %s\n", err.Error())
	}
	return
}

// JSON数组字符串反序列化成map数组
func JsonDecodeMapArray(str string) (thisMap []map[string]interface{}, err error) {
	if err := json.Unmarshal([]byte(str), &thisMap); err != nil {
		fmt.Printf("Json反序列化为Map出错: %s\n", err.Error())
	}
	return
}
