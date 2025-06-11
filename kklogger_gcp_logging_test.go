package kklogger_gcp_logging

import (
	"os"
	"testing"
	"time"

	"github.com/yetiz-org/goth-kklogger"
)

func TestKKLoggerRollbarHook(t *testing.T) {
	hook := &KKLoggerGCPLoggingHook{
		ProjectId:   os.Getenv("TEST_PROJECT_ID"),
		LogName:     "log_name",
		Environment: "production",
		CodeVersion: "code_version",
		Service:     "service",
		ServerRoot:  "server_root",
		Level:       kklogger.DebugLevel,
	}

	kklogger.AsyncWrite = false
	kklogger.SetLoggerHooks([]kklogger.LoggerHook{hook})
	kklogger.SetLogLevel("DEBUG")
	kklogger.TraceJ("tjsType", "jsData")
	kklogger.DebugJ("djsType", "jsData")
	kklogger.InfoJ("ijsType", "jsData")
	kklogger.WarnJ("wjsType", "jsData")
	kklogger.ErrorJ("ejsType", "jsData")
	time.Sleep(time.Second * 2)
}
