package cryptokit

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/spf13/cast"
	"github.com/text3cn/goodle/kit/timekit"
	"hash/crc32"
	"io/ioutil"
	"net/url"
	"strings"
)

// 生成MD5
func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	cipherStr := hash.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// 计算文件的 MD5 散列
func Md5File(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// 计算字符串的 SHA-1 散列
func Sha1(str string) string {
	hash := sha1.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

// 计算文件的 SHA-1 散列
func Sha1File(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	hash := sha1.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// 计算字符串的 32 位 CRC（循环冗余校验）
func Crc32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}

// 传统URL编码
func URLEncode(str string) string {
	return url.QueryEscape(str)
}

// 传统URL解码
func URLDecode(str string) (string, error) {
	return url.QueryUnescape(str)
}

// 遵守94年国际标准备忘录RFC 1738的URL编码，PHP中推荐使用，Rawurlencode，弃用URLEncode
func Rawurlencode(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

// Rawurldecode()的解码
func Rawurldecode(str string) (string, error) {
	return url.QueryUnescape(strings.Replace(str, "%20", "+", -1))
}

// 构造URL字符串
func HTTPBuildQuery(queryData url.Values) string {
	return queryData.Encode()
}

// Base64加密
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64解密
func Base64Decode(str string) (string, error) {
	switch len(str) % 4 {
	case 2:
		str += "=="
	case 3:
		str += "="
	}

	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func DynamicEncrypt(secret string, str string) string {
	return dynamicEncryption(secret, str, "encode")
}

func DynamicDecrypt(secret string, str string) string {
	return dynamicEncryption(secret, str, "decode")
}

// 可逆加解密算法, 不支持中文, Chr()转换问题，PHP的内置函数chr()支持中文
// 此方法从discuz的经典加解密算法移植而来，我封装的Chr()使用golang str[索引]来获取ASCII字符，所以无法支持中文
// 暂时没有去研究golang替代php::chr()函数的办法
// key : 加解密的密钥
// str ： 明文 / 密文
// operation ： encode / decode
// expire ：密文有效期，多少秒后失效，失效的密文解密后返回空，不传则不会过期
func dynamicEncryption(key string, str string, operation string, expire ...int) string {
	var expiry int
	if len(expire) > 0 {
		expiry = expire[0]
	}
	if str == "" {
		return "Token decryption invalid"
	}
	// 密文字符串中翻译掉浏览器不支持的特殊字符
	if operation == "decode" {
		str = strings.Replace(str, ".", "/", -1)
		str = strings.Replace(str, "-", "+", -1)
		str = strings.Replace(str, "_", "=", -1)
	}
	// 动态密匙长度，相同的明文会生成不同密文就是依靠动态密匙
	ckey_length := 4
	// 对传参进来的原始密匙进行MD5
	key = Md5(key)
	// 截取原始秘钥MD5的前16位然后再MD5一次作为秘钥a，用于参与加解密
	keya := Md5(key[0:16])
	// 截取原始秘钥MD5的后16位然后再MD5一次作为秘钥b，用来做数据完整性验证
	keyb := Md5(key[16:32])
	// 密匙c用于变化生成的密文, 解密时取要解密的密文的前四位
	var keyc string
	if operation == "encode" {
		keyc = Md5(cast.ToString(timekit.Microtime()))
		keyc = keyc[len(keyc)-ckey_length : len(keyc)]
	} else {
		keyc = str[0:ckey_length]
	}
	// 参与运算的秘钥
	cryptKey := keya + Md5(keya+keyc)
	// 参与运算的秘钥的长度
	key_length := len(cryptKey)
	// 解析str（明文或密文）
	if operation == "encode" {
		// 加密时
		expireTime := expiry
		if expiry > 0 {
			expireTime = expiry + timekit.NowTimestamp()
		}
		expireTimeString := fmt.Sprintf("%010d", expireTime)
		md5str := Md5(str + keyb)
		md5str = md5str[0:16] + str
		str = expireTimeString + md5str
	} else {
		// 解密时去掉密文的前 ckey_length 位，因为前ckey_length保存的是动态密码keyc
		str = str[ckey_length:]
		str, _ = Base64Decode(str)
	}
	string_length := len(str)
	result := ""
	// 生产一个数组，256个元素，每个数组元素的值等于数组key
	var box [127]int
	for i := 0; i < 127; i++ {
		box[i] = i
	}
	// 生产秘钥字典，从ASCII值获取字符串
	var randKey [127]uint8
	for i := 0; i <= 126; i++ {
		cryptKeyIndex := i % key_length
		randKey[i] = cryptKey[cryptKeyIndex]
	}
	// 用秘钥字典填充秘钥盒子 box
	for i, j := 0, 0; i < 127; i++ {
		j = (j + box[i] + int(randKey[i])) % 127
		tmp := box[i]
		box[i] = box[j]
		box[j] = tmp
	}
	// 最终加解密核心, 其实就是用这个算法打乱数组顺序，相当于是混淆
	// 因为异或运算符 c = a ^ b 可以推出 a =  b ^ c，所以加密时将位置交换过去，解密时交换回来，
	// 位运算符^是可逆加解密算法的核心
	for a, j, i := 0, 0, 0; i < string_length; i++ {
		a = (a + 1) % 127
		j = (j + box[a]) % 127
		tmp := box[a]
		box[a] = box[j]
		box[j] = tmp
		_i := int(str[i]) ^ (box[(box[a]+box[j])%127])
		result += Chr(_i)
	}
	// 加密得到密文
	if operation == "encode" {
		result = Base64Encode(result)
		result = strings.Replace(result, "/", ".", -1)
		result = strings.Replace(result, "+", "-", -1) // 中横线
		result = strings.Replace(result, "=", "_", -1) // 下划线
		result = keyc + result
		return result
	} else {
		// 解密，从result中截出明文,
		// 验证数据的有效性
		if result == "" {
			return ""
		}
		expireTime := cast.ToInt(result[0:10])
		resultFront10_16 := result[10:26] // dto : string
		verifyKey := Md5(result[26:] + keyb)
		verifyKey = verifyKey[0:16]
		// 1.验证是否过期；2.验证数据完整性
		if (expireTime == 0 || expireTime-timekit.NowTimestamp() > 0) && resultFront10_16 == verifyKey {
			return result[26:]
		}
		return ""
	}
}

// 将ASCII码值转化为字符串。
// 此函数与PHP的Mb_chr()函数转换结果一致，与php的chr()转换结果不一致
// 因为golang统一是utf-8编码，rune uses UTF-8，ASCII码值在127以下，127一下是可以和php对等，超过127的ASCII值翻译就无法对等了
func Chr(ascii int) string {
	return string(ascii)
}
