package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	//	. "smyhw.online/go/stdio_helper/helper"
	"smyhw.online/go/stdio_helper/helper/configer"
	"smyhw.online/go/stdio_helper/helper/logger"
)

var stdinPipe io.WriteCloser
var stdoutPipe io.ReadCloser

var opt_dir string = "./stdio_helper_dir"

var target_cmd *exec.Cmd

func main() {
	logger.Info("stdio helper启动...")
	//“检 查 文 件 完 整 性”
	logger.Info("校验目录完整性...")
	finfo, err := os.Stat(opt_dir)
	if os.IsNotExist(err) {
		logger.Info("操作文件夹不存在，创建...")
		err := os.Mkdir(opt_dir, 0777)
		if err != nil {
			logger.Warning("操作文件夹创建失败 -> " + err.Error())
			os.Exit(1)
		}
	} else if err != nil {
		logger.Warning("检查操作文件夹异常 -> " + err.Error())
		os.Exit(1)
	} else if !finfo.IsDir() {
		logger.Warning("当前目录下存在名为" + opt_dir + "的文件，stdio_helper无法启动")
		os.Exit(1)
	}
	//读取配置
	configer.Init()
	logger.Info("读取配置文件...")
	tmp1 := configer.GetString("target", "")
	if tmp1 == "" {
		logger.Warning("没有找到配置项目<target>,请检查配置文件...")
		os.Exit(1)
	}
	//启动目标程序
	logger.Info("启动目标程序<" + tmp1 + ">")
	//处理参数
	cd_arg := make([]string, len(os.Args)-1)
	for num, arg := range os.Args {
		if num == 0 {
			continue
		}
		cd_arg[num-1] = arg
	}
	target_cmd = exec.Command(tmp1, cd_arg...)
	stdinPipe, _ = target_cmd.StdinPipe()
	stdoutPipe, _ = target_cmd.StdoutPipe()
	//	target_cmd.Stdin = os.Stdin

	logger.Info("目标程序执行 -> " + target_cmd.String())
	err = target_cmd.Start()
	if err != nil {
		logger.Warning("启动目标程序失败 -> " + err.Error())
		os.Exit(1)
	}
	//设置running文件
	file, _ := os.Create(opt_dir + "/running")
	file.Close()
	logger.Info("进入主循环...")
	go handle_stdout()
	go handle_stdin()
	go handle_stdin_optfile()
	main_loop()
}

//获取目标程序stdout并打印到咱们自己的stdout上，一并写入opt_dir里的stdout文件里
func handle_stdout() {
	//	logger.Info("handle_stdout")
	//目标程序stdout的reader
	target_reader := bufio.NewReader(stdoutPipe)
	//打开std_out文件
	opt_file, err := os.OpenFile(opt_dir+"/stout", os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		logger.Warning("初始化stdout文件异常 -> " + err.Error())
		os.Exit(1)
	}
	for {
		tmp1, _, err := target_reader.ReadLine()
		outline := string(tmp1)

		if err != nil {
			if err.Error() == "EOF" {
				handle_target_exit()
			}
			logger.Warning("读取目标程序stdout异常" + err.Error())
		} else if outline != "" {
			fmt.Println(outline)
			opt_file.WriteString(outline + "\n")
		}
	}
}

//读取控制台输入
func handle_stdin() {
	reader := bufio.NewReader(os.Stdin)
	for {
		time.Sleep(0.1 * 1000000000)
		line, err := reader.ReadString('\n')
		if err != nil {
			logger.Warning("获取stdin出错 -> " + err.Error())
			continue
		}
		_, err = io.WriteString(stdinPipe, line)
		if err != nil {
			logger.Warning("转发stdin异常 -> " + err.Error())
			continue
		}
	}
}

//检测和读取stdin文件
func handle_stdin_optfile() {
	for {
		time.Sleep(0.1 * 1000000000)
		//		logger.Info("handle_opt_stdin")
		finfo, err := os.Stat(opt_dir + "/stdin")
		if os.IsNotExist(err) {
			continue
		}
		finfo.IsDir()
		file, err := os.Open(opt_dir + "/stdin")
		if err != nil {
			logger.Warning("检测stdin文件异常 -> " + err.Error())
			return
		}
		//读文件
		reader := bufio.NewReader(file)
		line, err := reader.ReadString('\n')
		for err != io.EOF {
			if line != "" {
				logger.Info("输入stdin -> " + line)
				io.WriteString(stdinPipe, line+"\n")
			}
			line, err = reader.ReadString('\n')
		}
		if line != "" {
			logger.Info("输入stdin -> " + line)
			io.WriteString(stdinPipe, line)
		}
		file.Close()
		os.Remove(opt_dir + "/stdin")
	}
}

//目标程序退出时调用
func handle_target_exit() {
	target_cmd.Wait()
	logger.Warning("目标程序退出 -> 退出码=" + fmt.Sprint(target_cmd.ProcessState.ExitCode()))
	err := os.Remove(opt_dir + "/running")
	if err != nil {
		logger.Warning("删除running文件失败 -> " + err.Error())
		os.Exit(1)
	}
	os.Exit(target_cmd.ProcessState.ExitCode())
}

//只是为了防止主线程退出
func main_loop() {
	for {
		time.Sleep(1000000000)
	}
}
