package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/Asutorufa/yuhaiin/config"
)

var (
	LockFilePath = config.Path + "/yuhaiin.lock"
	hostFile     = config.Path + "/host.txt"
	lockFile     *os.File
)

func GetProcessLock(str string) error {
	var err error
_retry:
	lockFile, err = os.OpenFile(LockFilePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path.Dir(LockFilePath), os.ModePerm)
			if err != nil {
				return fmt.Errorf("SettingEncodeJson():MkdirAll -> %v", err)
			}
			goto _retry
		}
		return fmt.Errorf("GetProcessLock() -> OpenFile() -> %v", err)
	}
	if err := LockFile(lockFile); err != nil {
		return fmt.Errorf("GetProcessLock() -> LockFile() -> %v", err)
	}
	err = ioutil.WriteFile(hostFile, []byte(str), os.ModePerm)
	if err != nil {
		log.Printf("GetProcessLock() -> WriteString() -> %v", err)
	}
	return nil
}

func ReadLockFile() (string, error) {
	s, err := ioutil.ReadFile(hostFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("ReadLockFile() -> ReadFile() -> %v", err)
	}
	return string(s), nil
}

func LockFileClose() (erra error) {
	err := os.Remove(hostFile)
	if err != nil {
		erra = fmt.Errorf("%v\nRemove hostFile -> %v", erra, err)
	}
	err = lockFile.Close()
	if err != nil {
		erra = fmt.Errorf("%v\nUnlock File (close file) -> %v", erra, err)
	}
	err = os.Remove(LockFilePath)
	if err != nil {
		erra = fmt.Errorf("%v\nRemove lockFile -> %v", erra, err)
	}
	return
}
