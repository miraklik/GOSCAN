package pkg

import (
	"log"
	"os"
	"os/user"
	"runtime"
)

type SysInfo struct {
	OS          string `json:"os"`
	ARCH        string `json:"arch"`
	Host        string `json:"hostname"`
	Username    string `json:"username"`
	UserHomeDir string `json:"home_dir"`
	Groupid     string `json:"gid"`
	Userid      string `json:"uid"`
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
		log.Printf("Warning: could not get user info: %v", err)
		return sysInfo
	}

	sysInfo.Username = u.Username
	sysInfo.UserHomeDir = u.HomeDir
	sysInfo.Groupid = u.Gid
	sysInfo.Userid = u.Uid

	return sysInfo
}
