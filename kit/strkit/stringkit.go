package strkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// 类型转字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func Tostring(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}

// 中横线拼接的字符串转驼峰式
func MethodNameToCamel(str string) string {
	// 切成数组
	stringArray := strings.Split(str, "-")
	// 首字母转大写
	for key, value := range stringArray {
		stringArray[key] = Ucfirst(value)
	}
	// 拼接
	returnString := strings.Join(stringArray, "")
	return returnString
}

// 生成唯一字符串，可以传入参数将唯一数值生成唯一字符串，比如手机号
// 用固定数值数值传按固定规则混淆，生成唯一字符串
func UniqueString(args ...string) string {
	var str string
	if len(args) == 0 {
		str = string(int(time.Now().Unix()))
	} else {
		str = args[0]
	}
	uniqueStr := Strtr(str, "1234567890", "huvtwkaemx")
	// 把顺序颠倒过来，然后再在中间找两个位置位置各插入一个字母
	var buffer bytes.Buffer
	for i := len(uniqueStr) - 1; i >= 0; i-- {
		buffer.WriteString(string(uniqueStr[i]))
		if i == 9 {
			buffer.WriteString(string(uniqueStr[3]))
		}
		if i == 5 {
			buffer.WriteString(string(uniqueStr[10]))
		}
	}
	uniqueStr = buffer.String()
	// 再交换一下顺序,藏头藏尾
	temp := Explode("", uniqueStr)
	_temp := temp[1]
	temp[1] = temp[7]
	temp[7] = _temp

	_temp = temp[2]
	temp[2] = temp[9]
	temp[9] = _temp

	_temp = temp[11]
	temp[11] = temp[5]
	temp[5] = _temp

	_temp = temp[12]
	temp[3] = temp[12]
	temp[12] = _temp

	var buffer2 bytes.Buffer
	for _, v := range temp {
		buffer2.WriteString(v)
	}
	uniqueStr = buffer.String()
	// 最后拼接个前缀
	uniqueStr = "mt" + uniqueStr
	return uniqueStr
}

// 生成唯一数字串，把时间戳的第一位砍掉，换成0-9的随机，第一位是1，到2033年，时间戳第一位变成2
func UniqueNumber() string {
	int64 := time.Now().Unix()
	timestamp := strconv.FormatInt(int64, 10)
	behead := timestamp[1:len(timestamp)]
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // 随机种子
	uniqueNumber := strconv.FormatInt(r.Int63n(8)+1, 10) + behead
	// 交换一下位置
	temp := Explode("", uniqueNumber)
	_temp := temp[1]
	temp[1] = temp[9]
	temp[9] = _temp

	_temp = temp[2]
	temp[2] = temp[7]
	temp[7] = _temp

	_temp = temp[4]
	temp[4] = temp[6]
	temp[6] = _temp

	var buffer bytes.Buffer
	for _, v := range temp {
		buffer.WriteString(v)
	}
	uniqueNumber = buffer.String()
	return uniqueNumber
}

// 生成随机字符串
func CreateNonceStr(length int) string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// 获取文件后缀
func GetSuffix(fileName string) string {
	s := Explode(".", fileName)
	last := s[len(s)-1]
	if last == fileName {
		return ""
	}
	return last
}

// 序列化成json
func JsonEncode(val interface{}) []byte {
	json, err := json.Marshal(val)
	if err != nil {
		fmt.Println("Json encode error")
	}
	return json
}

// JSON反序列化
func Json_decode_map(str string) (thisMap map[string]interface{}, err error) {
	if err = json.Unmarshal([]byte(str), &thisMap); err != nil {
		fmt.Printf("Json反序列化为Map出错: %s\n", err.Error())
	}
	return
}

// JSON数组字符串反序列化成map数组
func Json_decode_map_array(str string) (thisMap []map[string]interface{}, err error) {
	if err := json.Unmarshal([]byte(str), &thisMap); err != nil {
		fmt.Printf("Json反序列化为Map出错: %s\n", err.Error())
	}
	return
}

/**
 * 去除字符串首/尾逗号
 * @param str 要操作的字符串
 * @param mode 可选值：left 去除左边逗号，right 去除右边逗号。默认去除左右两边
 */
func TrimComma(str string, arg ...string) string {
	if str == "" || str == "," {
		return ""
	}
	mode := "ALL"
	if len(arg) > 0 {
		mode = arg[0]
	}
	str = Trim(str)
	ret := str
	switch strings.ToUpper(mode) {
	case "LEFT":
		if Substr(str, 0, 1) == "," {
			ret = str[1:]
		}
	case "RIGHT":
		if str[len(str)-1:] == "," {
			ret = str[0 : len(str)-1]
		}
	case "ALL":
		if Substr(str, 0, 1) == "," {
			ret = str[1:]
		}
		if ret[len(ret)-1:] == "," {
			ret = ret[0 : len(ret)-1]
		}
	}
	if ret == "" {
		ret = str
	}
	return ret
}

// 字符串切割成数组，返回字符串切片
func Explode(delimiter, str string) []string {
	if str == "" {
		return []string{}
	}
	return strings.Split(str, delimiter)
}

// 搜索字符串在另一字符串中是否存在，如果存在则返回该字符串及剩余部分，否则返回 FALSE。
func Strstr(haystack string, needle string) string {
	if needle == "" {
		return ""
	}
	idx := strings.Index(haystack, needle)
	if idx == -1 {
		return ""
	}
	return haystack[idx+len([]byte(needle))-1:]
}

// 字符串翻译函数，转换字符串中特定的字符。
// 如果params ...interface{}只传一个参数，类型是: map[string]string
// 例如：Strtr("baab", map[string]string{"ab": "01"}) 返回 "ba01"
// 如果params ...interface{}传两个参数, 类型是：string, string
// Strtr("baab", "ab", "01") 返回 "1001", a => 0; b => 1
func Strtr(haystack string, params ...interface{}) string {
	ac := len(params)
	if ac == 1 {
		pairs := params[0].(map[string]string)
		length := len(pairs)
		if length == 0 {
			return haystack
		}
		oldnew := make([]string, length*2)
		for o, n := range pairs {
			if o == "" {
				return haystack
			}
			oldnew = append(oldnew, o, n)
		}
		return strings.NewReplacer(oldnew...).Replace(haystack)
	} else if ac == 2 {
		from := params[0].(string)
		to := params[1].(string)
		trlen, lt := len(from), len(to)
		if trlen > lt {
			trlen = lt
		}
		if trlen == 0 {
			return haystack
		}

		str := make([]uint8, len(haystack))
		var xlat [256]uint8
		var i int
		var j uint8
		if trlen == 1 {
			for i = 0; i < len(haystack); i++ {
				if haystack[i] == from[0] {
					str[i] = to[0]
				} else {
					str[i] = haystack[i]
				}
			}
			return string(str)
		}
		// trlen != 1
		for {
			xlat[j] = j
			if j++; j == 0 {
				break
			}
		}
		for i = 0; i < trlen; i++ {
			xlat[from[i]] = to[i]
		}
		for i = 0; i < len(haystack); i++ {
			str[i] = xlat[haystack[i]]
		}
		return string(str)
	}

	return haystack
}

// 查找字符串在另一字符串中首次出现的位置（区分大小写）
func Strpos(haystack, needle string, offsetArg ...int) int {
	var offset int
	if len(offsetArg) == 0 {
		offset = 0
	} else {
		offset = offsetArg[0]
	}
	length := len(haystack)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	if offset < 0 {
		offset += length
	}
	pos := strings.Index(haystack[offset:], needle)
	if pos == -1 {
		return -1
	}
	return pos + offset
}

// 查找字符串在另一字符串中首次出现的位置（不区分大小写）
func Stripos(haystack, needle string, offset int) int {
	length := len(haystack)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	haystack = haystack[offset:]
	if offset < 0 {
		offset += length
	}
	pos := strings.Index(strings.ToLower(haystack), strings.ToLower(needle))
	if pos == -1 {
		return -1
	}
	return pos + offset
}

// 查找字符串在另一字符串中最后一次出现的位置（区分大小写）
// haystack : 被查找的字符串
// needle : 要在haystack中查找的字符串
// args[0] :  可选，规定从何处开始搜索
func Strrpos(haystack, needle string, args ...int) int {
	offset := 0
	if len(args) > 0 {
		offset = args[0]
	}
	pos, length := 0, len(haystack)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	if offset < 0 {
		haystack = haystack[:offset+length+1]
	} else {
		haystack = haystack[offset:]
	}
	pos = strings.LastIndex(haystack, needle)
	if offset > 0 && pos != -1 {
		pos += offset
	}
	return pos
}

// 查找字符串在另一字符串中最后一次出现的位置（不区分大小写）
func Strripos(haystack, needle string, offset int) int {
	pos, length := 0, len(haystack)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	if offset < 0 {
		haystack = haystack[:offset+length+1]
	} else {
		haystack = haystack[offset:]
	}
	pos = strings.LastIndex(strings.ToLower(haystack), strings.ToLower(needle))
	if offset > 0 && pos != -1 {
		pos += offset
	}
	return pos
}

// 字符串替换，在 subject 中将 old 替换成 new
func StrReplace(old, new, subject string, count ...int) string {
	num := -1
	if len(count) > 0 {
		num = count[0]
	}
	// -1 代表替换全部，0 代表不做替换，1 代表只替换一次
	return strings.Replace(subject, old, new, num)
}

// 字符串转大写
func Strtoupper(str string) string {
	return strings.ToUpper(str)
}

// 字符串转小写
func Strtolower(str string) string {
	return strings.ToLower(str)
}

// 字符串首字母转化为大写
func Ucfirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToUpper(v))
		return u + str[len(u):]
	}
	return ""
}

// 首字母转小写
func Lcfirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToLower(v))
		return u + str[len(u):]
	}
	return ""
}

// 单词首字母转大写
func Ucwords(str string) string {
	return strings.Title(str)
}

// 字符串截取
func Substr(str string, start uint, length int) string {
	if start < 0 || length < -1 {
		return str
	}
	switch {
	case length == -1:
		return str[start:]
	case length == 0:
		return ""
	}
	end := int(start) + length
	if end > len(str) {
		end = len(str)
	}
	return str[start:end]
}

// 字符串反转
func Strrev(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// 解析url查询字符串
// f1=m&f2=n -> map[f1:m f2:n]
// f[a]=m&f[b]=n -> map[f:map[a:m b:n]]
// f[a][a]=m&f[a][b]=n -> map[f:map[a:map[a:m b:n]]]
// f[]=m&f[]=n -> map[f:[m n]]
// f[a][]=m&f[a][]=n -> map[f:map[a:[m n]]]
// f[][]=m&f[][]=n -> map[f:[map[]]] // 不支持嵌套切片
// f=m&f[a]=n -> error // 这里和php不一样
// a .[[b=c -> map[a___[b:c]
func Parse_str(encodedString string, result map[string]interface{}) error {
	// build nested map.
	var build func(map[string]interface{}, []string, interface{}) error

	build = func(result map[string]interface{}, keys []string, value interface{}) error {
		length := len(keys)
		// trim ',"
		key := strings.Trim(keys[0], "'\"")
		if length == 1 {
			result[key] = value
			return nil
		}

		// The end is slice. like f[], f[a][]
		if keys[1] == "" && length == 2 {
			// todo nested slice
			if key == "" {
				return nil
			}
			val, ok := result[key]
			if !ok {
				result[key] = []interface{}{value}
				return nil
			}
			children, ok := val.([]interface{})
			if !ok {
				return fmt.Errorf("expected dto '[]interface{}' for key '%s', but got '%T'", key, val)
			}
			result[key] = append(children, value)
			return nil
		}

		// The end is slice + map. like f[][a]
		if keys[1] == "" && length > 2 && keys[2] != "" {
			val, ok := result[key]
			if !ok {
				result[key] = []interface{}{}
				val = result[key]
			}
			children, ok := val.([]interface{})
			if !ok {
				return fmt.Errorf("expected dto '[]interface{}' for key '%s', but got '%T'", key, val)
			}
			if l := len(children); l > 0 {
				if child, ok := children[l-1].(map[string]interface{}); ok {
					if _, ok := child[keys[2]]; !ok {
						_ = build(child, keys[2:], value)
						return nil
					}
				}
			}
			child := map[string]interface{}{}
			_ = build(child, keys[2:], value)
			result[key] = append(children, child)

			return nil
		}

		// map. like f[a], f[a][b]
		val, ok := result[key]
		if !ok {
			result[key] = map[string]interface{}{}
			val = result[key]
		}
		children, ok := val.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected dto 'map[string]interface{}' for key '%s', but got '%T'", key, val)
		}

		return build(children, keys[1:], value)
	}

	// split encodedString.
	parts := strings.Split(encodedString, "&")
	for _, part := range parts {
		pos := strings.Index(part, "=")
		if pos <= 0 {
			continue
		}
		key, err := url.QueryUnescape(part[:pos])
		if err != nil {
			return err
		}
		for key[0] == ' ' {
			key = key[1:]
		}
		if key == "" || key[0] == '[' {
			continue
		}
		value, err := url.QueryUnescape(part[pos+1:])
		if err != nil {
			return err
		}

		// split into multiple keys
		var keys []string
		left := 0
		for i, k := range key {
			if k == '[' && left == 0 {
				left = i
			} else if k == ']' {
				if left > 0 {
					if len(keys) == 0 {
						keys = append(keys, key[:left])
					}
					keys = append(keys, key[left+1:i])
					left = 0
					if i+1 < len(key) && key[i+1] != '[' {
						break
					}
				}
			}
		}
		if len(keys) == 0 {
			keys = append(keys, key)
		}
		// first key
		first := ""
		for i, chr := range keys[0] {
			if chr == ' ' || chr == '.' || chr == '[' {
				first += "_"
			} else {
				first += string(chr)
			}
			if chr == '[' {
				first += keys[0][i+1:]
				break
			}
		}
		keys[0] = first

		// build nested map
		if err := build(result, keys, value); err != nil {
			return err
		}
	}

	return nil
}

// 去除字符串两边空格
func Trim(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimSpace(str)
	}
	return strings.Trim(str, characterMask[0])
}

// 去除字符串左边空格
func Ltrim(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimLeftFunc(str, unicode.IsSpace)
	}
	return strings.TrimLeft(str, characterMask[0])
}

// 去除字符串右边空格
func Rtrim(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimRightFunc(str, unicode.IsSpace)
	}
	return strings.TrimRight(str, characterMask[0])
}

func ExplodeAndTrim(delimiter, str string) []string {
	arr := strings.Split(str, delimiter)
	var ret []string
	for _, v := range arr {
		ret = append(ret, Trim(v))
	}
	return ret
}

// string 转 int
func ParseInt(str string) int {
	data, _ := strconv.Atoi(str)
	return data
}

// string 转 int8
func ParseInt8(str string) int8 {
	data, _ := strconv.ParseInt(str, 10, 8)
	return int8(data)
}

// string 转 int32
func ParseInt32(str string) int32 {
	data, _ := strconv.ParseInt(str, 10, 32)
	return int32(data)
}

// string 转 int64
func ParseInt64(str string) int64 {
	// 如果是小数，
	data, _ := strconv.ParseInt(str, 10, 64)
	return data
}

func StringToFloat64(str string) float64 {
	data, _ := strconv.ParseFloat(str, 64)
	return data
}

// 判断字符串是否已 xxx 开头
func StartWith(str, prefix string) bool {
	return strings.HasPrefix(str, prefix)
}

// 判断字符串是否已 xxx 结尾
func EndWith(str, suffix string) bool {
	return strings.HasSuffix(str, suffix)
}

// 把字符串重复指定次数
func StrRepeat(input string, multiplier int) string {
	return strings.Repeat(input, multiplier)
}
