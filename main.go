package main

import (
	"github.com/yuyudeqiu/chronicle/cmd"
)

var (
	gitCommit  string
	gitDate    string
	buildTime  string
)

func main() {
	// 版本信息可以通过 chronicle version 查看
	cmd.Execute()
}
