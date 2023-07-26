package httpserver

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/spf13/cast"
	"io/ioutil"
	"mime/multipart"
)

// 为请求封装方法，在 Context 上实现接口
// 统一返回值中的 bool 代表请求方是否有传递这个数据过来
type IRequest interface {
	// 获取查询字符串中的参数，如: xxx.com?a=foo&b=bar&c[]=barbar
	Get(key string) interface{}
	GetInt(key string, defaultValue ...int) (int, bool)
	GetInt64(key string, defaultValue ...int64) (int64, bool)
	GetFloat64(key string, defaultValue ...float64) (float64, bool)
	GetFloat32(key string, defaultValue ...float32) (float32, bool)
	GetBool(key string, defaultValue ...bool) (bool, bool)
	GetString(key string, defaultValue ...string) (string, bool)
	GetStringSlice(key string, defaultValue ...[]string) ([]string, bool)

	// form 表单中的参数
	Form(key string) interface{}
	FormInt(key string, defaultValue ...int) (int, bool)
	FormInt64(key string, defaultValue ...int64) (int64, bool)
	FormFloat64(key string, defaultValue ...float64) (float64, bool)
	FormFloat32(key string, defaultValue ...float32) (float32, bool)
	FormBool(key string, defaultValue ...bool) (bool, bool)
	FormString(key string, defaultValue ...string) (string, bool)
	FormStringSlice(key string, defaultValue ...[]string) ([]string, bool)
	FormFile(key string) (*multipart.FileHeader, error)
	BindJson(obj interface{}) error // json body
	BindXml(obj interface{}) error  // xml body
	GetRawData() ([]byte, error)    // 其他格式

	// 其他格式
	Uri() string
	Method() string
	Host() string
	ClientIp() string

	// header
	Headers() map[string][]string
	Header(key string) (string, bool)

	// cookie
	Cookie(key string) (string, bool)
	Cookies() map[string]string
}

// 获取请求地址中所有参数
func (req *ReqStruct) QueryAll() map[string][]string {
	if req.request != nil {
		return req.request.URL.Query()
	}
	return map[string][]string{}
}

func (req *ReqStruct) Get(key string) interface{} {
	params := req.QueryAll()
	if vals, ok := params[key]; ok {
		return vals[0]
	}
	return nil
}

func (req *ReqStruct) GetInt(key string, defaultValue ...int) (int, bool) {
	params := req.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt(vals[0]), true
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0], false
	}
	return 0, false
}

func (req *ReqStruct) GetInt64(key string, defaultValue ...int64) (int64, bool) {
	params := req.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt64(vals[0]), true
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0], false
	}
	return 0, false
}

func (req *ReqStruct) GetFloat64(key string, defaultValue ...float64) (float64, bool) {
	params := req.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat64(vals[0]), true
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0], false
	}
	return 0, false
}

func (req *ReqStruct) GetFloat32(key string, defaultValue ...float32) (float32, bool) {
	params := req.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat32(vals[0]), true
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0], false
	}
	return 0, false
}

func (req *ReqStruct) GetBool(key string, defaultValue ...bool) (bool, bool) {
	params := req.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToBool(vals[0]), true
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0], false
	}
	return false, false
}

func (req *ReqStruct) GetString(key string, defaultValue ...string) (string, bool) {
	params := req.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return vals[0], true
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0], false
	}
	return "", false
}

// 请求 /xxx?a=11&a=22 中的参数 a 是能组成数组的
func (req *ReqStruct) GetStringSlice(key string, defaultValue ...[]string) ([]string, bool) {
	params := req.QueryAll()
	if vals, ok := params[key]; ok {
		return vals, true
	}
	if len(defaultValue) > 0 {
		return defaultValue[0], false
	}
	return []string{}, false
}

func (req *ReqStruct) FormAll() map[string][]string {
	if req.request != nil {
		req.request.ParseForm()
		return req.request.PostForm
	}
	return map[string][]string{}
}

func (req *ReqStruct) FormInt(key string, defaultValue ...int) (int, bool) {
	params := req.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt(vals[0]), true
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0], false
	}
	return 0, false
}

func (req *ReqStruct) FormInt64(key string, defaultValue ...int64) (int64, bool) {
	params := req.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt64(vals[0]), true
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0], false
	}
	return 0, false
}

func (req *ReqStruct) FormFloat64(key string, defaultValue ...float64) (float64, bool) {
	params := req.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat64(vals[0]), true
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0], false
	}
	return 0, false
}

func (req *ReqStruct) FormFloat32(key string, defaultValue ...float32) (float32, bool) {
	params := req.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat32(vals[0]), true
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0], false
	}
	return 0, false
}

func (req *ReqStruct) FormBool(key string, defaultValue ...bool) (bool, bool) {
	params := req.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToBool(vals[0]), true
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0], false
	}
	return false, false
}

func (req *ReqStruct) FormString(key string, defaultValue ...string) (string, bool) {
	params := req.FormAll()
	if vals, ok := params[key]; ok {
		return vals[0], true
	}
	if len(defaultValue) > 0 {
		return defaultValue[0], false
	}
	return "", false
}

func (req *ReqStruct) FormStringSlice(key string, defaultValue ...[]string) ([]string, bool) {
	params := req.FormAll()
	if vals, ok := params[key]; ok {
		return vals, true
	}
	if len(defaultValue) > 0 {
		return defaultValue[0], false
	}
	return []string{}, false
}

func (req *ReqStruct) FormFile(key string) (*multipart.FileHeader, error) {
	const defaultValueaultMultipartMemory = 32 << 20 // 32 MB
	if req.request.MultipartForm == nil {
		if err := req.request.ParseMultipartForm(defaultValueaultMultipartMemory); err != nil {
			return nil, err
		}
	}
	f, fh, err := req.request.FormFile(key)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, err
}

func (req *ReqStruct) Form(key string) interface{} {
	params := req.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return vals[0]
		}
	}
	return nil
}

// 将body文本解析到obj结构体中
func (req *ReqStruct) BindJson(obj interface{}) error {
	if req.request != nil {
		// 读取文本
		body, err := ioutil.ReadAll(req.request.Body)
		if err != nil {
			return err
		}
		// 重新填充request.Body，为后续的逻辑二次读取做准备
		req.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		// 解析到obj结构体中
		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("req.request empty")
	}
	return nil
}

// xml body
func (req *ReqStruct) BindXml(obj interface{}) error {
	if req.request != nil {
		body, err := ioutil.ReadAll(req.request.Body)
		if err != nil {
			return err
		}
		req.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		err = xml.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("req.request empty")
	}
	return nil
}

// 其他格式
func (req *ReqStruct) GetRawData() ([]byte, error) {
	if req.request != nil {
		body, err := ioutil.ReadAll(req.request.Body)
		if err != nil {
			return nil, err
		}
		req.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		return body, nil
	}
	return nil, errors.New("req.request empty")
}

// 基础信息
func (req *ReqStruct) Uri() string {
	return req.request.RequestURI
}

func (req *ReqStruct) Method() string {
	return req.request.Method
}

func (req *ReqStruct) Host() string {
	return req.request.URL.Host
}

func (req *ReqStruct) ClientIp() string {
	r := req.request
	ipAddress := r.Header.Get("X-Real-Ip")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}
	return ipAddress
}

// header
func (req *ReqStruct) Headers() map[string][]string {
	return req.request.Header
}

func (req *ReqStruct) Header(key string) (string, bool) {
	vals := req.request.Header.Values(key)
	if vals == nil || len(vals) <= 0 {
		return "", false
	}
	return vals[0], true
}

// cookie
func (req *ReqStruct) Cookies() map[string]string {
	cookies := req.request.Cookies()
	ret := map[string]string{}
	for _, cookie := range cookies {
		ret[cookie.Name] = cookie.Value
	}
	return ret
}

func (req *ReqStruct) Cookie(key string) (string, bool) {
	cookies := req.Cookies()
	if val, ok := cookies[key]; ok {
		return val, true
	}
	return "", false
}