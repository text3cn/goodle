package gokit

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
)

// 封装协程异常捕获，只管一层
// e.g: go SafeGo(fn)
// e.g: fn := func() { test("param") }; go SafeGo(fn)
func SafeGo(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[SafeGo panic recover]", err)
		}
	}()
	fn()
}

// 获取协程ID，慎用，仅用作有时候要打印协程ID仅当调试debug看一下用。
// 在C++中我们通过获取线程ID开辟不通空间保证线程安全。
// 而golang中google官方自从1.4就取消了获取协程ID的接口,不建议照C++那样做，因为滥用协程ID会导致GC无法及时回收内存。
// 说透彻点就是我知道你拿协程ID想干嘛,开数组/切片/map，然后根据协程ID来索引数据，这就会导致一个协程未退出你开的这个保存了一万个协程数据的大容器及其引用数据无法回收内存。
func GetGoroutineID() uint64 {
	b := make([]byte, 64)
	runtime.Stack(b, false)
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
