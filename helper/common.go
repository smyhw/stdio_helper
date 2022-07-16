package helper

import "smyhw.online/go/stdio_helper/helper/logger"

func FKnil(err error) {
	logger.Warning("发生异常 -> " + err.Error())
}
