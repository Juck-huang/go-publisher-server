package common

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
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
	cmd := exec.Command("/bin/bash", command...)
	var msg string
	if isPip {
		stdout, err := cmd.StdoutPipe()
		defer stdout.Close()
		if err != nil {
			msg = "执行命令失败"
			return msg, err
		}
		err = cmd.Start()
		if err != nil {
			msg = "执行命令失败"
			return msg, err
		}
		outBytes, err := ioutil.ReadAll(stdout)
		if err != nil {
			msg = "执行命令失败"
			return msg, err
		}
		splitBytes := strings.Split(string(outBytes), "\n")
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
