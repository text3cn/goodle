package cryptokit

//
///**
// * 注意：
// * golang 加解密常遇到的问题，与其他语言交互加解密时，凡是位运算的的需要注意选择具体的整型，例如：int16、int32
// * 之前实现 discuz 经典加解密算法时踩过坑，由于 ascii 字符编码不同 golang 中只能使用 int16 类型与 php 同步加解密
// * 这里同步 typescript 的加解密算法，不能使用 int 和 int64 来做位运算，又采坑一次。
// */
//import (
//	"math"
//	"statistics/kit/fn"
//	"strconv"
//
//	"statistics/kit/logger"
//
//	"github.com/gogf/gf/errors/gerror"
//)
//
//// 是否已经通过密钥初始化过转译字典，这个计算在程序的整个生命周期只需要计算一次，放到 boot 中去
//var isInit = false
//
//// 转译字典，通过传入密钥动态生成
//var cipherMap, numMap map[string]string
//
//var psdOrigins = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
//
//const HEAD_BYTE = 0xB0
//
//func _init(key int) {
//	if isInit {
//		// 使用全局变量做单例
//		return
//	}
//	isInit = true
//	getMap(key)
//	// numMap 是 cipherMap 的 key 和 value 对调
//	numMap = make(map[string]string)
//	for k, v := range cipherMap {
//		numMap[v] = k
//	}
//}
//
//func getMap(key int) {
//	// 十六进制的16个数
//	table := []string{}
//	for {
//		if len(psdOrigins) == 0 {
//			break
//		}
//		time := math.Ceil(float64(len(psdOrigins)) / float64(key))
//		ids := []int{}
//		rangeDesc := rangeRight(int(time))
//		for _, v := range rangeDesc {
//			index := key * fn.ParseInt(string(v))
//			table = append(table, psdOrigins[index])
//			ids = append(ids, index)
//		}
//		// 从 psdOrigins 中删除 ids
//		for _, v := range ids {
//			psdOrigins = fn.ArrayStringDeleteOne(v, psdOrigins)
//		}
//	}
//	// 使用数组填充 map，用 table.index 做 map 的 key（十六进制），用 table.value 做 map 的 value
//	cipherMap = make(map[string]string)
//	for k, v := range table {
//		hex := strconv.FormatInt(int64(k), 16)
//		cipherMap[fn.Tostring(hex)] = v
//	}
//	return
//}
//
//// 生成从数字几开始(不包括本身)，到 0 为止的连续数字串
//func rangeRight(start int) (ret []byte) {
//	str := ""
//	for i := start - 1; i >= 0; i-- {
//		str += fn.Tostring(i)
//	}
//	ret = []byte(str)
//	return
//}
//
///*
//位移操作：(int32(value) >> (int32(8) * it)) & 0xFF 结果与 typescript 不一致，加密失败。
//有符号与无符号的整型都无法与 ts 操作位移达成一致，原因未深究。
//界面函数可以用，本项目目前只需要界面，不需要加密，遂放弃之。
//对应 ts 的加密方法在：lion:src/kit/util/IdEncryptionKit.ts 中
//func EncrypNumber(value, key int) (cipher string, err error) {
//	_init(key)
//	rangeNums := rangeRight(8)
//	bytes := []int32{}
//	for _, v := range rangeNums {
//		it := int32(fn.ParseInt(string(v)))
//		item := (int32(value) >> (int32(8) * it)) & 0xFF
//		bytes = append(bytes, item)
//	}
//	checkSum := checkSum(bytes)
//	idBytes := []int32{}
//	idBytes = append(idBytes, HEAD_BYTE)
//	for _, v := range bytes {
//		idBytes = append(idBytes, v)
//	}
//	idBytes = append(idBytes, checkSum)
//	// 获得密文
//	var prev []int32 = []int32{}
//	var idNumbers []int32 = []int32{}
//	reduceFn := func(prev []int32, it int32) []int32 {
//		ret := []int32{}
//		if it >= 16 {
//			ret = append(ret, int32(it>>4))
//			ret = append(ret, int32(it&0xF))
//		} else {
//			ret = append(ret, 0)
//			ret = append(ret, it)
//		}
//		return ret
//	}
//	for _, it := range idBytes {
//		prev = idNumbers
//		idNumbers = append(idNumbers, reduceFn(prev, it)...)
//	}
//	logger.Cred(idNumbers)
//	for i, it := range idNumbers {
//		index := encodeIndex(it, int32(i))
//		s := cipherMap[fn.Tostring(index)]
//		hex := fn.Dechex(fn.ParseInt64(s))
//		cipher += hex
//	}
//	return
//}*/
//
//// 解密 typescript 加密的 userId，key 固定为 3 ，两边保持一致。
//func DecrypNovelUserId(cipher string) string {
//	uid, err := DecrypNumber(cipher, 3)
//	if err != nil {
//		logger.Error(err)
//		return ""
//	}
//	str := strconv.Itoa(int(uid))
//	return str
//}
//
//func DecrypNumber(cipher string, key int) (userId int32, err error) {
//	// 初始化转义字典, key 为密钥
//	_init(key)
//	if len(cipher) != 20 {
//		err = gerror.New("密文长度不对")
//	}
//	cipherByte := []byte(cipher)
//	charsNumberArr := []int{}
//	for i := 0; i < 20; i++ {
//		charsNumberArr = append(charsNumberArr)
//	}
//	chars := []int32{}
//	for i, it := range cipherByte {
//		hex := numMap[string(it)]
//		dec, err := fn.Hexdec(hex)
//		if err != nil {
//			err = gerror.New("进制转换失败")
//		}
//		chars = append(chars, int32(decodeIndex(int(dec), i)))
//	}
//	chunkArr := fn.ChunkForIntArry(chars, 2)
//	numbers := []int32{}
//	for _, it := range chunkArr {
//		num := (it[0] << 4) + it[1]
//		numbers = append(numbers, num)
//	}
//	head := numbers[0]
//	tail := numbers[len(numbers)-1]
//	bytes := numbers[1 : len(numbers)-1]
//	checkSum := checkSum(bytes)
//	if head != HEAD_BYTE || checkSum != tail {
//		err = gerror.New("密文" + cipher + "不正确")
//	}
//	// 获得明文
//	var prev int32 = 0
//	var value int32 = 0
//	reduceFn := func(prev, it int32) int32 {
//		bit := prev << 8
//		ret := bit + it
//		return ret
//	}
//	for _, it := range bytes {
//		prev = value
//		value = reduceFn(prev, int32(it))
//	}
//	userId = int32(value)
//	return userId, err
//}
//
//func decodeIndex(value, index int) int32 {
//	return int32((32 + value - index) % 16)
//}
//
//func encodeIndex(value, index int32) int32 {
//	return (value + index) % 16
//}
//
//func checkSum(bytes []int32) int32 {
//	var num int32 = 0
//	for _, v := range bytes {
//		num += v
//	}
//	return num & 0xFF
//}

//func ChunkForIntArry(s []int32, size int) [][]int32 {
//	if size < 1 {
//		panic("size 不能小于 1")
//	}
//	length := len(s)
//	chunks := int32(math.Ceil(float64(length) / float64(size)))
//	var n [][]int32
//	for i, end := 0, 0; chunks > 0; chunks-- {
//		end = (i + 1) * size
//		if end > length {
//			end = length
//		}
//		n = append(n, s[i*size:end])
//		i++
//	}
//	return n
//}
