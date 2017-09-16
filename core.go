package main

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
	"strings"

	"github.com/donpark/pam"
)

const logPrefix = "PAMTLS "

var logger *log.Logger

func coreInit(args pam.Args) {
	if loggerSpec, loggerSpecified := args["logger"]; loggerSpecified {
		if strings.ToLower(loggerSpec) == "syslog" {
			initSyslog()
		} else {
			logger = log.New(os.Stdout, logPrefix, 0)
			logger.Printf("Unrecognised logger spec: %s", loggerSpec)
		}
	} else {
		logger = log.New(os.Stdout, logPrefix, log.Ltime)
	}
}

func initSyslog() {
	var err error
	logger, err = syslog.NewLogger(syslog.LOG_INFO, log.Ltime)
	if err != nil {
		logger = log.New(os.Stderr, logPrefix, log.Ltime)
		logger.Printf("syslog.Open() err: %v", err)
	}
}

func info(module string, data ...interface{}) {
	s := fmt.Sprintf("[%s] %s", module, fmt.Sprint(data...))
	logger.Println(s)
}
