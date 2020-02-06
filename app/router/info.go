package router

import (
	"github.com/gin-gonic/gin"
	"github.com/goecology/webhook/app/pkg/conf"
	"os/exec"
)

var (
	CODEMAP = map[string]string{
		"exit status 0":     "命令成功结束",
		"exit status 1":     "通用未知错误",
		"exit status 2":     "误用Shell命令",
		"exit status 126":   "命令不可执行",
		"exit status 127":   "没找到命令",
		"exit status 128":   "无效的退出参数",
		"exit status 128+x": "Linux信号x的严重错误",
		"exit status 130":   "命令通过Ctrl+C控制码越界",
		"exit status 255":   "退出码越界",
	}
)

func Info(c *gin.Context) {
	for _, value := range conf.Conf.WebHook {
		if value.UrlPath == c.Request.URL.Path {
			token := c.Request.URL.Query().Get("token")
			if value.Token != "" && value.Token != token {
				c.String(200, "%s", "token is error")
				return
			}

			outputline := execCommand(value.IsBash, value.ExecPath, value.ExecParams)
			c.String(200, "%s", outputline)
			return
		}
	}

}

func execCommand(isBash bool, execPath string, execParams []string) string {
	if isBash {
		return runBashCommand(execPath, execParams)
	}
	return runCommonCommand(execPath, execParams)
}

func runBashCommand(execPath string, execParams []string) (output string) {
	args := make([]string, 0)
	args = append(args, "-c")
	args = append(args, execPath)
	for _, value := range execParams {
		args = append(args, value)
	}
	command := exec.Command("/bin/bash", args...)
	resp, err := command.CombinedOutput()
	if err != nil {
		errStr := errMap(err.Error())
		output = "runBashCommand info: " + errStr + ",  error: " + err.Error()
		return
	}
	output = string(resp)
	return
}

func runCommonCommand(execPath string, execParams []string) (output string) {
	command := exec.Command(execPath, execParams...)
	resp, err := command.CombinedOutput()
	if err != nil {
		errStr := errMap(err.Error())
		output = "runCommonCommand info: " + errStr + ",  error: " + err.Error()
		return
	}
	output = string(resp)
	return
}

func errMap(code string) string {
	errStr, flag := CODEMAP[code]
	if !flag {
		return ""
	}
	return errStr
}
