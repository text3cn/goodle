package httpserver

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

// 为响应封装方法
// 方法返回 IResponse 本身代表支持链式调用，例如：res.SetOkStatus().Json("success")
type IResponse interface {
	Json(obj interface{}) IResponse
	Html(template string, obj interface{}) IResponse
	Jsonp(obj interface{}) IResponse
	Xml(obj interface{}) IResponse
	Text(format string, values ...interface{}) IResponse
	Redirect(path string) IResponse // 重定向
	SetHeader(key string, val string) IResponse
	SetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse
	SetOkStatus() IResponse       // 设置 200 状态
	SetStatus(code int) IResponse // 设置其他状态码
}

func (res *RespStruct) Json(obj interface{}) IResponse {
	byt, err := json.Marshal(obj)
	if err != nil {
		return res.SetStatus(http.StatusInternalServerError)
	}
	res.SetHeader("Content-Type", "application/json")
	res.responseWriter.Write(byt)
	return res
}

// Jsonp输出
func (res *RespStruct) Jsonp(obj interface{}) IResponse {
	// 获取请求参数callback
	callbackFunc := res.request.GetString("callback", "callback_function")
	res.SetHeader("Content-Type", "application/javascript")
	// 输出到前端页面的时候需要注意下进行字符过滤，否则有可能造成xss攻击
	callback := template.JSEscapeString(callbackFunc)

	// 输出函数名
	_, err := res.responseWriter.Write([]byte(callback))
	if err != nil {
		return res
	}
	// 输出左括号
	_, err = res.responseWriter.Write([]byte("("))
	if err != nil {
		return res
	}
	// 数据函数参数
	ret, err := json.Marshal(obj)
	if err != nil {
		return res
	}
	_, err = res.responseWriter.Write(ret)
	if err != nil {
		return res
	}
	// 输出右括号
	_, err = res.responseWriter.Write([]byte(")"))
	if err != nil {
		return res
	}
	return res
}

// xml输出
func (res *RespStruct) Xml(obj interface{}) IResponse {
	byt, err := xml.Marshal(obj)
	if err != nil {
		return res.SetStatus(http.StatusInternalServerError)
	}
	res.SetHeader("Content-Type", "application/html")
	res.responseWriter.Write(byt)
	return res
}

// html输出
func (res *RespStruct) Html(file string, obj interface{}) IResponse {
	// 读取模版文件，创建template实例
	t, err := template.New("output").ParseFiles(file)
	if err != nil {
		return res
	}
	// 执行Execute方法将obj和模版进行结合
	if err := t.Execute(res.responseWriter, obj); err != nil {
		return res
	}

	res.SetHeader("Content-Type", "text/html")
	return res
}

// string
func (res *RespStruct) Text(format string, values ...interface{}) IResponse {
	out := fmt.Sprintf(format, values...)
	res.SetHeader("Content-Type", "text/plain")
	res.responseWriter.Write([]byte(out))
	return res
}

// 重定向
func (res *RespStruct) Redirect(path string) IResponse {
	http.Redirect(res.responseWriter, res.request.request, path, http.StatusMovedPermanently)
	return res
}

// header
func (res *RespStruct) SetHeader(key string, val string) IResponse {
	res.responseWriter.Header().Add(key, val)
	return res
}

// Cookie
func (res *RespStruct) SetCookie(key string, val string, maxAge int, path string, domain string, secure bool, httpOnly bool) IResponse {
	if path == "" {
		path = "/"
	}
	http.SetCookie(res.responseWriter, &http.Cookie{
		Name:     key,
		Value:    url.QueryEscape(val),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: 1,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
	return res
}

// 设置状态码
func (res *RespStruct) SetStatus(code int) IResponse {
	res.responseWriter.WriteHeader(code)
	return res
}

// 设置200状态
func (res *RespStruct) SetOkStatus() IResponse {
	res.responseWriter.WriteHeader(http.StatusOK)
	return res
}
