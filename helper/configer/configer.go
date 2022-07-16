package configer

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"smyhw.online/go/stdio_helper/helper/logger"
)

type StdioHelperCfg struct {
	Target string `json:"target"`
}

var cfg StdioHelperCfg

func fknil(e error) {
	if e != nil {
		logger.Warning("读取配置文件异常")
		panic(e)
	}
}

func Init() {
	cfg_file_str := "./stdio_helper_dir/config.json"
	logger.Info("init config...")
	//检测配置文件
	finfo, err := os.Stat(cfg_file_str)
	if os.IsNotExist(err) {
		logger.Info("配置文件不存在，创建...")
		file, err := os.Create(cfg_file_str)
		fknil(err)
		file.WriteString("{\"target\": \"./awa.exe\"}")
		file.Close()
	} else if err != nil {
		fknil(err)
	} else {
		if finfo.IsDir() {
			logger.Warning("配置文件变文件夹了？？？ -> " + cfg_file_str)
			os.Exit(1)
		}
	}

	//读取配置文件
	cfg_byte, err := ioutil.ReadFile(cfg_file_str)
	fknil(err)
	ee := json.Unmarshal(cfg_byte, &cfg)
	fknil(ee)
	//	logger.Info(fmt.Sprint("读取到配置文件 --> ", cfg))
}

func get_raw(key string) string {
	//TODO 完成它...
	return cfg.Target
}

func GetString(key string, def string) string {
	re := get_raw(key)
	if re == "" {
		return def
	}
	return re
}
