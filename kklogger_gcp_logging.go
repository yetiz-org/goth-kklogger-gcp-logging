package kklogger_gcp_logging

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"cloud.google.com/go/logging"
	kklogger "github.com/yetiz-org/goth-kklogger"
)

type KKLoggerGCPLoggingHook struct {
	enabled     bool
	initOnce    sync.Once
	logger      *logging.Logger
	ProjectId   string
	LogName     string
	Environment string
	CodeVersion string
	Service     string
	ServerRoot  string
	Level       kklogger.Level
}

func (h *KKLoggerGCPLoggingHook) LogString(args ...interface{}) string {
	if args == nil {
		return ""
	}

	args = args[0].([]interface{})
	argl := len(args)

	if argl == 1 {
		switch tp := args[0].(type) {
		case string:
			return tp
		}
	} else if argl > 1 {
		switch tp := args[0].(type) {
		case string:
			pargs := args[1:]
			return fmt.Sprintf(tp, pargs...)
		}
	}

	return fmt.Sprint(args...)
}

func (h *KKLoggerGCPLoggingHook) Trace(args ...interface{}) {
	if h.Level < kklogger.TraceLevel {
		return
	}

	h.Send(kklogger.TraceLevel, h.LogString(args...))
}

func (h *KKLoggerGCPLoggingHook) Debug(args ...interface{}) {
	if h.Level < kklogger.DebugLevel {
		return
	}

	h.Send(kklogger.DebugLevel, h.LogString(args...))
}

func (h *KKLoggerGCPLoggingHook) Info(args ...interface{}) {
	if h.Level < kklogger.InfoLevel {
		return
	}

	h.Send(kklogger.InfoLevel, h.LogString(args...))
}

func (h *KKLoggerGCPLoggingHook) Warn(args ...interface{}) {
	if h.Level < kklogger.WarnLevel {
		return
	}

	h.Send(kklogger.WarnLevel, h.LogString(args...))
}

func (h *KKLoggerGCPLoggingHook) Error(args ...interface{}) {
	if h.Level < kklogger.ErrorLevel {
		return
	}

	h.Send(kklogger.ErrorLevel, h.LogString(args...))
}

func (h *KKLoggerGCPLoggingHook) Send(level kklogger.Level, msg string) {
	h.initOnce.Do(func() {
		client, err := logging.NewClient(context.Background(), h.ProjectId)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		h.logger = client.Logger(h.LogName)
		h.enabled = true
	})

	if !h.enabled {
		return
	}

	h.logger.Log(h.getEntry(level, msg))
}

func (h *KKLoggerGCPLoggingHook) getEntry(level kklogger.Level, msg string) logging.Entry {
	obj := &map[string]interface{}{}
	json.Unmarshal([]byte(msg), obj)
	return logging.Entry{
		Severity: func(level kklogger.Level) logging.Severity {
			switch level {
			case kklogger.TraceLevel:
				return logging.Default
			case kklogger.DebugLevel:
				return logging.Debug
			case kklogger.InfoLevel:
				return logging.Info
			case kklogger.WarnLevel:
				return logging.Warning
			case kklogger.ErrorLevel:
				return logging.Error
			}

			return logging.Default
		}(level),
		Payload: obj,
		Labels: map[string]string{
			"logName":     h.LogName,
			"environment": h.Environment,
			"codeVersion": h.CodeVersion,
			"service":     h.Service,
			"serverRoot":  h.ServerRoot,
		},
	}
}
