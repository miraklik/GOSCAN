package pkg

import (
	"log"
	"os"
	"os/user"
	"runtime"
)

type SysInfo struct {
	OS          string
	ARCH        string
	Host        string
	Groupid     string
	Userid      string
	Username    string
	UserHomeDir string
}

func GetSysInfo() SysInfo {
	var sysInfo SysInfo

	sysInfo.OS = runtime.GOOS
	sysInfo.ARCH = runtime.GOARCH

	name, err := os.Hostname()
	if err == nil {
		sysInfo.Host = name
	}

	u, err := user.Current()
	if err != nil {
		log.Fatalf("user.Current() failed: %v", err)
	}
	sysInfo.Username = u.Username
	sysInfo.UserHomeDir = u.HomeDir
	sysInfo.Groupid = u.Gid
	sysInfo.Userid = u.Uid

	return sysInfo
}
