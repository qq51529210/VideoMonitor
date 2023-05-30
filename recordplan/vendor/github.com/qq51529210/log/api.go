package log

import "os"

var (
	// 默认
	defaultLogger *Logger
	// 级别
	debugLevel = "[D] "
	infoLevel  = "[I] "
	warnLevel  = "[W] "
	errorLevel = "[E] "
	panicLevel = "[P] "
)

func init() {
	// 默认 logger
	defaultLogger = NewLogger(os.Stdout, new(DefaultHeader), "")
}

// SetLogger 设置默认的 Logger
func SetLogger(lg *Logger) {
	defaultLogger = lg
}

// GetLogger 返回默认的 Logger
func GetLogger() *Logger {
	return defaultLogger
}

// Recover 使用默认 Logger 函数
func Recover(recover any) {
	defaultLogger.Recover(recover)
}

// IsDebug 返回是否启用 debug
func IsDebug() bool {
	return !defaultLogger.disableDebug
}

// EnableDebug 设置是否启用 debug
func EnableDebug(enable bool) {
	defaultLogger.disableDebug = !enable
}

// Debug 输出日志
func Debug(args ...any) {
	if defaultLogger.disableDebug {
		return
	}
	defaultLogger.print(loggerDepth, &debugLevel, args...)
}

// Debugf 输出日志
func Debugf(format string, args ...any) {
	if defaultLogger.disableDebug {
		return
	}
	defaultLogger.printf(loggerDepth, &debugLevel, &format, args...)
}

// DebugDepth 输出日志
func DebugDepth(depth int, args ...any) {
	if defaultLogger.disableDebug {
		return
	}
	defaultLogger.print(loggerDepth+depth, &debugLevel, args...)
}

// DebugfDepth 输出日志
func DebugfDepth(depth int, format string, args ...any) {
	if defaultLogger.disableDebug {
		return
	}
	defaultLogger.printf(loggerDepth+depth, &debugLevel, &format, args...)
}

// DebugTrace 输出日志
func DebugTrace(traceID string, args ...any) {
	if defaultLogger.disableDebug {
		return
	}
	if traceID != "" {
		defaultLogger.printTrace(loggerDepth, &traceID, &debugLevel, args...)
	} else {
		defaultLogger.print(loggerDepth, &debugLevel, args...)
	}
}

// DebugfTrace 输出日志
func DebugfTrace(traceID, format string, args ...any) {
	if defaultLogger.disableDebug {
		return
	}
	if traceID != "" {
		defaultLogger.printfTrace(loggerDepth, &traceID, &debugLevel, &format, args...)
	} else {
		defaultLogger.printf(loggerDepth, &debugLevel, &format, args...)
	}
}

// DebugDepthTrace 输出日志
func DebugDepthTrace(depth int, traceID string, args ...any) {
	if defaultLogger.disableDebug {
		return
	}
	if traceID != "" {
		defaultLogger.printTrace(loggerDepth+depth, &traceID, &debugLevel, args...)
	} else {
		defaultLogger.print(loggerDepth+depth, &debugLevel, args...)
	}
}

// DebugfDepthTrace 输出日志
func DebugfDepthTrace(depth int, traceID, format string, args ...any) {
	if defaultLogger.disableDebug {
		return
	}
	if traceID != "" {
		defaultLogger.printfTrace(loggerDepth+depth, &traceID, &debugLevel, &format, args...)
	} else {
		defaultLogger.printf(loggerDepth+depth, &debugLevel, &format, args...)
	}
}

// IsInfo 返回是否启用 info
func IsInfo() bool {
	return !defaultLogger.disableInfo
}

// EnableInfo 设置是否启用 info
func EnableInfo(enable bool) {
	defaultLogger.disableInfo = !enable
}

// Info 输出日志
func Info(args ...any) {
	if defaultLogger.disableInfo {
		return
	}
	defaultLogger.print(loggerDepth, &infoLevel, args...)
}

// Infof 输出日志
func Infof(format string, args ...any) {
	if defaultLogger.disableInfo {
		return
	}
	defaultLogger.printf(loggerDepth, &infoLevel, &format, args...)
}

// InfoDepth 输出日志
func InfoDepth(depth int, args ...any) {
	if defaultLogger.disableInfo {
		return
	}
	defaultLogger.print(loggerDepth+depth, &infoLevel, args...)
}

// InfofDepth 输出日志
func InfofDepth(depth int, format string, args ...any) {
	if defaultLogger.disableInfo {
		return
	}
	defaultLogger.printf(loggerDepth+depth, &infoLevel, &format, args...)
}

// InfoTrace 输出日志
func InfoTrace(traceID string, args ...any) {
	if defaultLogger.disableInfo {
		return
	}
	if traceID != "" {
		defaultLogger.printTrace(loggerDepth, &traceID, &infoLevel, args...)
	} else {
		defaultLogger.print(loggerDepth, &infoLevel, args...)
	}
}

// InfofTrace 输出日志
func InfofTrace(traceID, format string, args ...any) {
	if defaultLogger.disableInfo {
		return
	}
	if traceID != "" {
		defaultLogger.printfTrace(loggerDepth, &traceID, &infoLevel, &format, args...)
	} else {
		defaultLogger.printf(loggerDepth, &infoLevel, &format, args...)
	}
}

// InfoDepthTrace 输出日志
func InfoDepthTrace(depth int, traceID string, args ...any) {
	if defaultLogger.disableInfo {
		return
	}
	if traceID != "" {
		defaultLogger.printTrace(loggerDepth+depth, &traceID, &infoLevel, args...)
	} else {
		defaultLogger.print(loggerDepth+depth, &infoLevel, args...)
	}
}

// InfofDepthTrace 输出日志
func InfofDepthTrace(depth int, traceID, format string, args ...any) {
	if defaultLogger.disableInfo {
		return
	}
	if traceID != "" {
		defaultLogger.printfTrace(loggerDepth+depth, &traceID, &infoLevel, &format, args...)
	} else {
		defaultLogger.printf(loggerDepth+depth, &infoLevel, &format, args...)
	}
}

// IsWarn 返回是否启用 warn
func IsWarn() bool {
	return !defaultLogger.disableWarn
}

// EnableWarn 设置是否启用 warn
func EnableWarn(enable bool) {
	defaultLogger.disableWarn = !enable
}

// Warn 输出日志
func Warn(args ...any) {
	if defaultLogger.disableWarn {
		return
	}
	defaultLogger.print(loggerDepth, &warnLevel, args...)
}

// Warnf 输出日志
func Warnf(format string, args ...any) {
	if defaultLogger.disableWarn {
		return
	}
	defaultLogger.printf(loggerDepth, &warnLevel, &format, args...)
}

// WarnDepth 输出日志
func WarnDepth(depth int, args ...any) {
	if defaultLogger.disableWarn {
		return
	}
	defaultLogger.print(loggerDepth+depth, &warnLevel, args...)
}

// WarnfDepth 输出日志
func WarnfDepth(depth int, format string, args ...any) {
	if defaultLogger.disableWarn {
		return
	}
	defaultLogger.printf(loggerDepth+depth, &warnLevel, &format, args...)
}

// WarnTrace 输出日志
func WarnTrace(traceID string, args ...any) {
	if defaultLogger.disableWarn {
		return
	}
	if traceID != "" {
		defaultLogger.printTrace(loggerDepth, &traceID, &warnLevel, args...)
	} else {
		defaultLogger.print(loggerDepth, &warnLevel, args...)
	}
}

// WarnfTrace 输出日志
func WarnfTrace(traceID, format string, args ...any) {
	if defaultLogger.disableWarn {
		return
	}
	if traceID != "" {
		defaultLogger.printfTrace(loggerDepth, &traceID, &warnLevel, &format, args...)
	} else {
		defaultLogger.printf(loggerDepth, &warnLevel, &format, args...)
	}
}

// WarnDepthTrace 输出日志
func WarnDepthTrace(depth int, traceID string, args ...any) {
	if defaultLogger.disableWarn {
		return
	}
	if traceID != "" {
		defaultLogger.printTrace(loggerDepth+depth, &traceID, &warnLevel, args...)
	} else {
		defaultLogger.print(loggerDepth+depth, &warnLevel, args...)
	}
}

// WarnfDepthTrace 输出日志
func WarnfDepthTrace(depth int, traceID, format string, args ...any) {
	if defaultLogger.disableWarn {
		return
	}
	if traceID != "" {
		defaultLogger.printfTrace(loggerDepth+depth, &traceID, &warnLevel, &format, args...)
	} else {
		defaultLogger.printf(loggerDepth+depth, &warnLevel, &format, args...)
	}
}

// IsError 返回是否启用 error
func IsError() bool {
	return !defaultLogger.disableError
}

// EnableError 设置是否启用 error
func EnableError(enable bool) {
	defaultLogger.disableError = !enable
}

// Error 输出日志
func Error(args ...any) {
	if defaultLogger.disableError {
		return
	}
	defaultLogger.print(loggerDepth, &errorLevel, args...)
}

// Errorf 输出日志
func Errorf(format string, args ...any) {
	if defaultLogger.disableError {
		return
	}
	defaultLogger.printf(loggerDepth, &errorLevel, &format, args...)
}

// ErrorDepth 输出日志
func ErrorDepth(depth int, args ...any) {
	if defaultLogger.disableError {
		return
	}
	defaultLogger.print(loggerDepth+depth, &errorLevel, args...)
}

// ErrorfDepth 输出日志
func ErrorfDepth(depth int, format string, args ...any) {
	if defaultLogger.disableError {
		return
	}
	defaultLogger.printf(loggerDepth+depth, &errorLevel, &format, args...)
}

// ErrorTrace 输出日志
func ErrorTrace(traceID string, args ...any) {
	if defaultLogger.disableError {
		return
	}
	if traceID != "" {
		defaultLogger.printTrace(loggerDepth, &traceID, &errorLevel, args...)
	} else {
		defaultLogger.print(loggerDepth, &errorLevel, args...)
	}
}

// ErrorfTrace 输出日志
func ErrorfTrace(traceID, format string, args ...any) {
	if defaultLogger.disableError {
		return
	}
	if traceID != "" {
		defaultLogger.printfTrace(loggerDepth, &traceID, &errorLevel, &format, args...)
	} else {
		defaultLogger.printf(loggerDepth, &errorLevel, &format, args...)
	}
}

// ErrorDepthTrace 输出日志
func ErrorDepthTrace(depth int, traceID string, args ...any) {
	if defaultLogger.disableError {
		return
	}
	if traceID != "" {
		defaultLogger.printTrace(loggerDepth+depth, &traceID, &errorLevel, args...)
	} else {
		defaultLogger.print(loggerDepth+depth, &errorLevel, args...)
	}
}

// ErrorfDepthTrace 输出日志
func ErrorfDepthTrace(depth int, traceID, format string, args ...any) {
	if defaultLogger.disableError {
		return
	}
	if traceID != "" {
		defaultLogger.printfTrace(loggerDepth+depth, &traceID, &errorLevel, &format, args...)
	} else {
		defaultLogger.printf(loggerDepth+depth, &errorLevel, &format, args...)
	}
}
