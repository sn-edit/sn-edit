package conf

import (
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strconv"
)

func Err(err interface{}, fields log.Fields, exit bool) {

	// add the original caller from runtime
	_, fn, line, _ := runtime.Caller(1)
	log.WithFields(log.Fields{"caller_fn": fn + ":" + strconv.Itoa(line)}).WithFields(fields).Error(err)

	if exit {
		os.Exit(1)
	}
}
