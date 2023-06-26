package common

import (
	"errors"
	"hy.juck.com/go-publisher-server/config"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
)

var (
	G = config.G
)

// EleIsExistSlice 判断元素是否存在数组中
func EleIsExistSlice(element string, sliceElement []string) bool {
	index := sort.SearchStrings(sliceElement, element)
	if index < len(sliceElement) && sliceElement[index] == element {
		return true
	}
	return false
}

// ExecCommand 执行命令,分管道和标准输出
func ExecCommand(isPip bool, command ...string) (string, error) {
	G.Logger.Infof("执行的发布命令为：%s", command)
	cmd := exec.Command("/bin/bash", command...)
	var msg string
	if isPip {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			G.Logger.Errorf("执行命令失败：失败原因[%s]", err.Error())
			msg = "执行命令失败"
			return msg, err
		}
		err = cmd.Start()
		if err != nil {
			G.Logger.Errorf("执行命令失败：失败原因[%s]", err.Error())
			msg = "执行命令失败"
			return msg, err
		}
		outBytes, err := ioutil.ReadAll(stdout)
		if err != nil {
			G.Logger.Errorf("执行命令失败：失败原因[%s]", err.Error())
			msg = "执行命令失败"
			return msg, err
		}
		defer stdout.Close()
		splitBytes := strings.Split(string(outBytes), "\n")
		G.Logger.Infof("执行发布命令结果:%s", splitBytes)
		if len(splitBytes) > 1 {
			return splitBytes[len(splitBytes)-2], nil
		} else {
			return "", errors.New("发布失败")
		}

	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return "", err
		}
		return "", nil
	}
}
