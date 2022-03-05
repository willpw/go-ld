package LD

import (
	"github.com/ying32/govcl/vcl"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	Dir, _  = os.Getwd()
	Ldir    = ""
	Cmdi    = 2
	Apkname = ""
	Ldid    sync.Map // Ldid 存储脚本状态
)

// 雷电list2命令返回数据处理返回窗口句柄
func list2(byt []byte, i string) string {
	st := string(byt)
	st2 := strings.Split(st, "\r\n")
	id := 0
	for _, val := range st2 {
		st3 := strings.Split(val, ",")
		//0-ID 1-模拟器名称 2-模拟器顶层窗口句柄 3-模拟器绑定窗口句柄 4-是否进入android（0|1） 5-进程PID 6-VBox进程PID
		if st3[0] == i && st3[4] == "1" {
			return st3[3]
		}
		if st3[0] == i {
			//判断模拟器是否存在
			id = 1
		}
	}
	if id == 0 {
		runtime.Goexit()
	}
	return "0"
}

// 		雷电CMD命令集合
// 		i-- 选择CMD执行方式 key-- CMD参数
func Ldcmd(cmdi int, key ...string) string {

	if file(Ldir, "dnconsole.exe") == false {
		//判断模拟器目录是否
		vcl.ShowMessage("模拟器目录错误！")
		runtime.Goexit()
	}
	cmd := exec.Command("", "")
	cm := []string{Ldir + "dnconsole.exe", key[0]}
	switch key[0] {
	case "launch", "isrunning", "reboot", "quit":
		cm = append(cm, "--index", key[1])
	case "runapp", "killapp", "launchex":
		cm = append(cm, "--index", key[1], "--packagename", Apkname)
	case "modify":
		cm = append(cm, "--index", key[1])
		cm = append(cm, "--resolution", "540,960,240")
		cm = append(cm, "--cpu", "2")
		cm = append(cm, "--memory", "1024")
		cm = append(cm, "--lockwindow", "1")
		cm = append(cm, "--autorotate", "0")
	case "list2", "sortWnd", "quitall":
	default:
		vcl.ShowMessage("参数错误！！")
		runtime.Goexit()
	}
	cmd.Path = cm[0]
	cmd.Args = append(cm)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	switch cmdi {
	case 1:
		byt, err := cmd.CombinedOutput()
		if err != nil {
			if len(byt) > 1500 {
				vcl.ShowMessage("cmd命令执行----失败！-1\n" + string(byt[:1500]))
				return ""
			}
			vcl.ShowMessage("cmd命令执行----失败！-1\n" + string(byt))
			return ""
		}
		if len(byt) > 1500 {
			vcl.ShowMessage("cmd命令执行----成功！-1\n" + string(byt[:1500]))
		}
		vcl.ShowMessage("cmd命令执行----成功！-1\n" + string(byt))

		switch key[0] {
		case "list2":
			return list2(byt, key[1])
		case "isrunning":
			return string(byt)
		}
	case 2:
		byt, err := cmd.Output()
		if err != nil {
			vcl.ShowMessage("cmd命令执行----失败！-2")
			return ""
		}
		switch key[0] {
		case "list2":
			return list2(byt, key[1])
		case "isrunning":
			return string(byt)
		}
	case 3:
		err := cmd.Wait()
		if err != nil {
			vcl.ShowMessage("cmd命令执行----失败！-3")
		}
	case 4:
		err := cmd.Start()
		if err != nil {
			vcl.ShowMessage("cmd命令执行----失败！-4")
		}
	case 5:
		err := cmd.Run()
		if err != nil {
			vcl.ShowMessage("cmd命令执行----失败！-5")
		}
	}
	return ""
}

//模拟器检测启动并运行APP
func ldrun(i string) {
	if Ldcmd(Cmdi, "isrunning", i) == "running" {
		Ldcmd(Cmdi, "killapp", i)
		time.Sleep(time.Millisecond * 1000)
		Ldcmd(Cmdi, "runapp", i)
		return
	}
	Ldcmd(Cmdi, "list2", i)  //识别模拟器是否存在
	Ldcmd(Cmdi, "modify", i) //对模拟器进行设置
	Ldcmd(Cmdi, "launchex", i)
	for n := 0; Ldcmd(Cmdi, "isrunning", i) != "running"; n++ {
		if n > 10 {
			runtime.Goexit()
		}
		time.Sleep(time.Millisecond * 15000)
	}

}

//雷电按键输入
func Ldanj(key ...string) {

}

// 检测文件是否存在
func file(dir, pat string) bool {
	_, err := os.Stat(dir + pat)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {

		return false
	}
	vcl.ShowMessage("文件存在检测错误！！\n" + pat)
	return false
}
