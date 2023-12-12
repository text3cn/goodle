package syskit

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime/debug"
)

// 执行系统命令，并且返回执行的命令和执行结果
func ExecSysCmd(commandName string, params []string) (string, string) {
	cmd := exec.Command(commandName, params...)
	//fmt.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	strings := ""
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		//fmt.Println(line)
		strings += line
	}
	cmd.Wait()
	commandStr := ""
	for _, v := range cmd.Args {
		commandStr += " " + v
	}
	return commandStr, strings
}

// 获取当前进程 ID
func Getpid() int {
	return os.Getpid()
}

// 获取调用栈
func GetStack() string {
	// 打印： debug.PrintStack()
	return string(debug.Stack())
}
