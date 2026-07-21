package main

import (
	"github.com/yuyudeqiu/chronicle/cmd"
	"github.com/yuyudeqiu/chronicle/internal/buildinfo"
)

// Build variables injected via ldflags (make build) or runtime/debug.ReadBuildInfo (go install)
var (
	version   string
	gitCommit string
	gitDate   string
	buildTime string
)

func main() {
	buildinfo.Set(version, gitCommit, gitDate, buildTime)
	cmd.Execute()
}
