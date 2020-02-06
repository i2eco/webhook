package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goecology/webhook/app/pkg/conf"
	"io"
	"os/exec"
	"strings"
)

//
//var (
//	CODEMAP = map[string]string{
//		"exit status 0":     "命令成功结束",
//		"exit status 1":     "通用未知错误",
//		"exit status 2":     "误用Shell命令",
//		"exit status 126":   "命令不可执行",
//		"exit status 127":   "没找到命令",
//		"exit status 128":   "无效的退出参数",
//		"exit status 128+x": "Linux信号x的严重错误",
//		"exit status 130":   "命令通过Ctrl+C控制码越界",
//		"exit status 255":   "退出码越界",
//	}
//)

func Info(c *gin.Context) {
	for _, value := range conf.Conf.WebHook {
		fmt.Println("value.UrlPath ==>", value.UrlPath)
		if value.UrlPath == c.Request.URL.Path {
			token := c.Request.URL.Query().Get("token")
			if value.Token != "" && value.Token != token {
				c.String(200, "%s", "token is error")
				return
			}
			execCommand(c, value.IsBash, value.ExecPath, value.ExecParams)
			return
		}
	}
	c.String(200, "%s", "no webhook url")
	return

}

func execCommand(c *gin.Context, isBash bool, execPath string, execParams []string) {
	if isBash {
		runBashCommand(c, execPath, execParams)
		return
	}
	runCommonCommand(c, execPath, execParams)
}

//
//func errMap(code string) string {
//	errStr, flag := CODEMAP[code]
//	if !flag {
//		return ""
//	}
//	return errStr
//}

func runBashCommand(c *gin.Context, execPath string, execParams []string) {
	args := make([]string, 0)
	args = append(args, "-c")
	args = append(args, execPath)
	for _, value := range execParams {
		args = append(args, value)
	}
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", args...)

	//显示运行的命令
	c.String(200, "%s", cmd.Args)
	//

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		c.String(200, "Error starting command: %s......", err.Error())
		return
	}

	go asyncLog(c, stdout)
	go asyncLog(c, stderr)

	if err := cmd.Wait(); err != nil {
		c.String(200, "Error waiting for command execution: %s......", err.Error())
		return
	}
	return
}

func runCommonCommand(c *gin.Context, execPath string, execParams []string) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command(execPath, execParams...)

	//显示运行的命令
	c.String(200, "%s", cmd.Args)

	//显示运行的命令
	c.String(200, "%s", cmd.Args)
	//

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		c.String(200, "Error starting command: %s......", err.Error())
		return
	}

	go asyncLog(c, stdout)
	go asyncLog(c, stderr)

	if err := cmd.Wait(); err != nil {
		c.String(200, "Error waiting for command execution: %s......", err.Error())
		return
	}
	return
}

func asyncLog(c *gin.Context, reader io.ReadCloser) error {
	cache := "" //缓存不足一行的日志信息
	buf := make([]byte, 1024)
	for {
		num, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if num > 0 {
			b := buf[:num]
			s := strings.Split(string(b), "\n")
			line := strings.Join(s[:len(s)-1], "\n") //取出整行的日志
			c.String(200, "%s%s\n", cache, line)
			cache = s[len(s)-1]
		}
	}
}
