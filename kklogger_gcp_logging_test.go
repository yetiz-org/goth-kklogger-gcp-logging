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
	kklogger.ReportCaller = true
	kklogger.SetLoggerHooks([]kklogger.LoggerHook{hook})
	kklogger.SetLogLevel("DEBUG")
	kklogger.TraceJ("tjsType", "jsData")
	kklogger.DebugJ("djsType", "jsData")
	kklogger.InfoJ("ijsType", "jsData")
	kklogger.WarnJ("wjsType", "jsData")
	kklogger.ErrorJ("ejsType", "jsData")
	time.Sleep(time.Second * 2)
}

func TestExtendedLoggerHookImplementation(t *testing.T) {
	hook := &KKLoggerGCPLoggingHook{
		ProjectId:   "test-project",
		LogName:     "test-log",
		Environment: "test",
		CodeVersion: "v1.0.0",
		Service:     "test-service",
		ServerRoot:  "/test",
		Level:       kklogger.InfoLevel,
	}

	var extHook kklogger.ExtendedLoggerHook = hook
	if extHook == nil {
		t.Fatal("Hook does not implement ExtendedLoggerHook interface")
	}

	testFunc := "test.function"
	testFile := "/path/to/file.go"
	testLine := 42

	hook.InfoWithCaller(testFunc, testFile, testLine, "test message")
	hook.ErrorWithCaller(testFunc, testFile, testLine, "error message")
}

func TestBasicHookBackwardCompatibility(t *testing.T) {
	hook := &KKLoggerGCPLoggingHook{
		ProjectId:   "test-project",
		LogName:     "test-log",
		Environment: "test",
		CodeVersion: "v1.0.0",
		Service:     "test-service",
		ServerRoot:  "/test",
		Level:       kklogger.InfoLevel,
	}

	var basicHook kklogger.LoggerHook = hook
	if basicHook == nil {
		t.Fatal("Hook does not implement LoggerHook interface")
	}

	basicHook.Info("test message")
	basicHook.Error("error message")
}
