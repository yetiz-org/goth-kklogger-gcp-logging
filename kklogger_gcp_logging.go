package kklogger_gcp_logging

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
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
	if args == nil || len(args) == 0 {
		return ""
	}

	if len(args) == 1 {
		if slice, ok := args[0].([]interface{}); ok {
			args = slice
		}
	}

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

	h.Send(kklogger.TraceLevel, "", "", 0, h.LogString(args...))
}

func (h *KKLoggerGCPLoggingHook) Debug(args ...interface{}) {
	if h.Level < kklogger.DebugLevel {
		return
	}

	h.Send(kklogger.DebugLevel, "", "", 0, h.LogString(args...))
}

func (h *KKLoggerGCPLoggingHook) Info(args ...interface{}) {
	if h.Level < kklogger.InfoLevel {
		return
	}

	h.Send(kklogger.InfoLevel, "", "", 0, h.LogString(args...))
}

func (h *KKLoggerGCPLoggingHook) Warn(args ...interface{}) {
	if h.Level < kklogger.WarnLevel {
		return
	}

	h.Send(kklogger.WarnLevel, "", "", 0, h.LogString(args...))
}

func (h *KKLoggerGCPLoggingHook) Error(args ...interface{}) {
	if h.Level < kklogger.ErrorLevel {
		return
	}

	h.Send(kklogger.ErrorLevel, "", "", 0, h.LogString(args...))
}

func (h *KKLoggerGCPLoggingHook) TraceWithCaller(funcName, file string, line int, args ...interface{}) {
	if h.Level < kklogger.TraceLevel {
		return
	}

	h.Send(kklogger.TraceLevel, funcName, file, line, h.LogString(args...))
}

func (h *KKLoggerGCPLoggingHook) DebugWithCaller(funcName, file string, line int, args ...interface{}) {
	if h.Level < kklogger.DebugLevel {
		return
	}

	h.Send(kklogger.DebugLevel, funcName, file, line, h.LogString(args...))
}

func (h *KKLoggerGCPLoggingHook) InfoWithCaller(funcName, file string, line int, args ...interface{}) {
	if h.Level < kklogger.InfoLevel {
		return
	}

	h.Send(kklogger.InfoLevel, funcName, file, line, h.LogString(args...))
}

func (h *KKLoggerGCPLoggingHook) WarnWithCaller(funcName, file string, line int, args ...interface{}) {
	if h.Level < kklogger.WarnLevel {
		return
	}

	h.Send(kklogger.WarnLevel, funcName, file, line, h.LogString(args...))
}

func (h *KKLoggerGCPLoggingHook) ErrorWithCaller(funcName, file string, line int, args ...interface{}) {
	if h.Level < kklogger.ErrorLevel {
		return
	}

	h.Send(kklogger.ErrorLevel, funcName, file, line, h.LogString(args...))
}

func (h *KKLoggerGCPLoggingHook) Send(level kklogger.Level, funcName, file string, line int, msg string) {
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

	h.logger.Log(h.getEntry(level, funcName, file, line, msg))
}

func (h *KKLoggerGCPLoggingHook) getEntry(level kklogger.Level, funcName, file string, line int, msg string) logging.Entry {
	obj := map[string]interface{}{}
	_ = json.Unmarshal([]byte(msg), &obj)
	labels := map[string]string{
		"logProject":  h.LogName,
		"environment": h.Environment,
		"codeVersion": h.CodeVersion,
		"service":     h.Service,
		"serverRoot":  h.ServerRoot,
	}

	if funcName != "" {
		labels["caller_function"] = funcName
	}
	if file != "" {
		labels["caller_file"] = file
	}
	if line > 0 {
		labels["caller_line"] = fmt.Sprintf("%d", line)
	}

	if v, f := obj["type"]; f {
		if typeStr, ok := v.(string); ok {
			// Parse format: current_file_package_name:struct_name.method_name#section_name!action_tag
			// Split by ':' to get package name
			if colonIdx := strings.Index(typeStr, ":"); colonIdx != -1 {
				labels["log_package"] = typeStr[:colonIdx]
				remaining := typeStr[colonIdx+1:]

				// Split by '#' to separate method from section/action
				if hashIdx := strings.Index(remaining, "#"); hashIdx != -1 {
					classMethod := remaining[:hashIdx]
					sectionAction := remaining[hashIdx+1:]

					// Split struct_name.method_name
					if dotIdx := strings.Index(classMethod, "."); dotIdx != -1 {
						labels["log_class"] = classMethod[:dotIdx]
						labels["log_method"] = classMethod[dotIdx+1:]
					} else {
						labels["log_method"] = classMethod
					}

					// Split by '!' to separate section from action
					if exclamIdx := strings.Index(sectionAction, "!"); exclamIdx != -1 {
						labels["log_section"] = sectionAction[:exclamIdx]
						labels["log_action"] = sectionAction[exclamIdx+1:]
					} else {
						labels["log_section"] = sectionAction
					}
				} else {
					// No section/action part, check if there's action only
					if exclamIdx := strings.Index(remaining, "!"); exclamIdx != -1 {
						classMethod := remaining[:exclamIdx]
						labels["log_action"] = remaining[exclamIdx+1:]

						// Split struct_name.method_name
						if dotIdx := strings.Index(classMethod, "."); dotIdx != -1 {
							labels["log_class"] = classMethod[:dotIdx]
							labels["log_method"] = classMethod[dotIdx+1:]
						} else {
							labels["log_method"] = classMethod
						}
					} else {
						// Split struct_name.method_name
						if dotIdx := strings.Index(remaining, "."); dotIdx != -1 {
							labels["log_class"] = remaining[:dotIdx]
							labels["log_method"] = remaining[dotIdx+1:]
						} else {
							labels["log_method"] = remaining
						}
					}
				}
			}
		}
	}

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
		Payload: &obj,
		Labels:  labels,
	}
}
