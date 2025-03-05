package util

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
)

const (
	SUDO_UID = "SUDO_UID"
	SUDO_GID = "SUDO_GID"
)

func GetEnvInt(key string) (parsed int, err error) {
	if value := os.Getenv(key); value == "" {
		err = fmt.Errorf("missing env key %s", key)
	} else if parsed, err = strconv.Atoi(value); err != nil {
		err = fmt.Errorf("cannot parse env key %s=%q: %v", key, value, err)
	}
	return
}

var sudo_uid int
var sudo_gid int

func CheckPrivilege() {
	var err error
	euid := os.Geteuid()
	sudo_uid, err = GetEnvInt(SUDO_UID)
	if euid != 0 || err != nil {
		fmt.Fprintln(os.Stderr, "should run with sudo")
		os.Exit(1)
	}
	egid := os.Getegid()
	sudo_gid, err = GetEnvInt(SUDO_GID)
	if egid != 0 || err != nil {
		fmt.Fprintln(os.Stderr, "should run with sudo")
		os.Exit(1)
	}
}

func GetPrivilege() {
	err := syscall.Seteuid(0)
	if err != nil {
		fmt.Fprintln(os.Stderr, "set euid 0 error", err)
		os.Exit(1)
	}
	err = syscall.Setegid(0)
	if err != nil {
		fmt.Fprintln(os.Stderr, "set egid 0 error", err)
		os.Exit(1)
	}
}

func DropPrivilege(temporarily bool) {
	if temporarily {
		syscall.Seteuid(sudo_uid)
		syscall.Setegid(sudo_gid)
	} else {
		syscall.Setuid(sudo_uid)
		syscall.Setgid(sudo_gid)
	}
}
