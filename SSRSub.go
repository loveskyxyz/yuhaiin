// +build !noGui

package main

import (
	"github.com/Asutorufa/SsrMicroClient/config"
	"github.com/Asutorufa/SsrMicroClient/gui"
	"github.com/Asutorufa/SsrMicroClient/init"
	"github.com/Asutorufa/SsrMicroClient/process/lockfile"
	"log"
	//_ "net/http/pprof"
	"os"
)

var (
	lockFile = config.GetConfigAndSQLPath() + "/SsrMicroClientRunStatuesLockFile"
)

func main() {
	//go func() {
	//	// 开启pprof，监听请求
	//	ip := "0.0.0.0:6060"
	//	if err := http.ListenAndServe(ip, nil); err != nil {
	//		fmt.Printf("start pprof failed on %s\n", ip)
	//	}
	//}()
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	configPath := config.GetConfigAndSQLPath()
	ssrinit.Init(configPath)
	lockFile, err := os.Create(lockFile)
	if err != nil {
		log.Println(err)
		return
	}
	if err = lockfile.LockFile(lockFile); err != nil {
		log.Println("process is exist!\n" + err.Error())
		return
	}
	defer func() {
		_ = lockFile.Close()
		_ = os.Remove(configPath + "/SsrMicroClientRunStatuesLockFile")
	}()

	ssrMicroClientGUI, err := gui.NewSsrMicroClientGUI()
	if err != nil {
		if ssrMicroClientGUI != nil {
			ssrMicroClientGUI.MessageBox(err.Error())
		} else {
			log.Println(err)
		}
	}
	if ssrMicroClientGUI != nil {
		//ssrMicroClientGUI.MainWindow.Show()
		ssrMicroClientGUI.App.Exec()
	} else {
		log.Println("gui is nil")
	}
}
