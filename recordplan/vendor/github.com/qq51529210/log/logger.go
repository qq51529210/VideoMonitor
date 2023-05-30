package log

import (
	"fmt"
	"io"
	"runtime"
)

var (
	// header:
	headerEnd = []byte(": ")
)

const (
	loggerDepth = 3
)

// Logger 默认实现
type Logger struct {
	// 输出
	io.Writer
	// 头格式
	Header
	// 名称
	name string
	// 是否禁止 debug
	disableDebug bool
	// 是否禁止 info
	disableInfo bool
	// 是否禁止 warn
	disableWarn bool
	// 是否禁止 error
	disableError bool
}

// NewLogger 返回默认的 Logger
// 格式 "Header [name] [level] [tracID] text"
func NewLogger(writer io.Writer, header Header, name string) *Logger {
	lg := new(Logger)
	lg.Writer = writer
	lg.Header = header
	if name != "" {
		lg.name = fmt.Sprintf("[%s] ", name)
	}
	return lg
}

func (lg *Logger) print(depth int, level *string, args ...any) {
	l := logPool.Get().(*Log)
	l.b = l.b[:0]
	// 名称
	if lg.name != "" {
		l.b = append(l.b, lg.name...)
	}
	// 级别
	l.b = append(l.b, *level...)
	// 头
	lg.Header.Time(l)
	lg.Header.Stack(l, depth)
	l.b = append(l.b, headerEnd...)
	// 日志
	fmt.Fprint(l, args...)
	// 换行
	l.b = append(l.b, '\n')
	// 输出
	lg.Writer.Write(l.b)
	// 回收
	logPool.Put(l)
}

func (lg *Logger) printf(depth int, level, format *string, args ...any) {
	l := logPool.Get().(*Log)
	l.b = l.b[:0]
	// 名称
	if lg.name != "" {
		l.b = append(l.b, lg.name...)
	}
	// 级别
	l.b = append(l.b, *level...)
	// 头
	lg.Header.Time(l)
	l.b = append(l.b, ' ')
	lg.Header.Stack(l, depth)
	l.b = append(l.b, headerEnd...)
	// 日志
	fmt.Fprintf(l, *format, args...)
	// 换行
	l.b = append(l.b, '\n')
	// 输出
	lg.Writer.Write(l.b)
	// 回收
	logPool.Put(l)
}

func (lg *Logger) printTrace(depth int, trace, level *string, args ...any) {
	l := logPool.Get().(*Log)
	l.b = l.b[:0]
	// 名称
	if lg.name != "" {
		l.b = append(l.b, lg.name...)
	}
	// 级别
	l.b = append(l.b, *level...)
	// 头
	lg.Header.Time(l)
	l.b = append(l.b, ' ')
	lg.Header.Stack(l, depth)
	l.b = append(l.b, headerEnd...)
	// 追踪
	l.b = append(l.b, '[')
	l.b = append(l.b, *trace...)
	l.b = append(l.b, ']')
	l.b = append(l.b, ' ')
	// 日志
	fmt.Fprint(l, args...)
	// 换行
	l.b = append(l.b, '\n')
	// 输出
	lg.Writer.Write(l.b)
	// 回收
	logPool.Put(l)
}

func (lg *Logger) printfTrace(depth int, trace, level, format *string, args ...any) {
	l := logPool.Get().(*Log)
	l.b = l.b[:0]
	// 名称
	if lg.name != "" {
		l.b = append(l.b, lg.name...)
	}
	// 级别
	l.b = append(l.b, *level...)
	// 头
	lg.Header.Time(l)
	l.b = append(l.b, ' ')
	lg.Header.Stack(l, depth)
	l.b = append(l.b, headerEnd...)
	// 追踪
	l.b = append(l.b, '[')
	l.b = append(l.b, *trace...)
	l.b = append(l.b, ']')
	l.b = append(l.b, ' ')
	// 日志
	fmt.Fprintf(l, *format, args...)
	// 换行
	l.b = append(l.b, '\n')
	// 输出
	lg.Writer.Write(l.b)
	// 回收
	logPool.Put(l)
}

// Recover 如果 recover 不为 nil，输出堆栈
func (lg *Logger) Recover(recover any) {
	if recover == nil {
		return
	}
	// get statck info l.line
	b := logPool.Get().(*Log)
	b.b = b.b[:cap(b.b)]
	for {
		n := runtime.Stack(b.b, false)
		if n < len(b.b) {
			b.b = b.b[:n]
			break
		}
		b.b = make([]byte, len(b.b)+1024)
	}
	// log
	l := logPool.Get().(*Log)
	l.b = l.b[:0]
	// 名称
	if lg.name != "" {
		l.b = append(l.b, lg.name...)
	}
	// 级别
	l.b = append(l.b, panicLevel...)
	lg.Header.Time(l)
	// recover
	fmt.Fprintf(l, ": %v\n", recover)
	// filter
	p := b.b
	// 找到 panic.go
	found := false
	for len(p) > 0 {
		// find new line
		i := 0
		for ; i < len(p); i++ {
			if p[i] == '\n' {
				i++
				break
			}
		}
		line := p[:i]
		p = p[i:]
		// find file line
		if line[0] != '\t' {
			continue
		}
		if !found {
			found = hasPanicGO(line)
			continue
		}
		// \t filepath/file.go:line +0x622
		i = len(line) - 1
		for i > 0 {
			if line[i] == ' ' {
				//
				line = line[:i]
				break
			}
			i--
		}
		// write
		l.b = append(l.b, "[statck] "...)
		l.b = append(l.b, line[1:]...)
		l.b = append(l.b, '\n')
	}
	// 输出
	lg.Writer.Write(l.b)
	// 回收
	logPool.Put(b)
	logPool.Put(l)
}

func hasPanicGO(line []byte) bool {
	for i := len(line) - 1; i > 1; i-- {
		if line[i] == '/' {
			for j := i; j < len(line); j++ {
				if line[j] == 'p' &&
					line[j+1] == 'a' &&
					line[j+2] == 'n' &&
					line[j+3] == 'i' &&
					line[j+4] == 'c' &&
					line[j+5] == '.' &&
					line[j+6] == 'g' &&
					line[j+7] == 'o' {
					return true
				}
			}
			return false
		}
	}
	return false
}

// IsDebug 返回是否启用 debug
func (lg *Logger) IsDebug() bool {
	return !lg.disableDebug
}

// EnableDebug 设置是否启用 debug
func (lg *Logger) EnableDebug(enable bool) {
	lg.disableDebug = !enable
}

// Debug 输出日志
func (lg *Logger) Debug(args ...any) {
	if lg.disableDebug {
		return
	}
	lg.print(loggerDepth, &debugLevel, args...)
}

// Debugf 输出日志
func (lg *Logger) Debugf(format string, args ...any) {
	if lg.disableDebug {
		return
	}
	lg.printf(loggerDepth, &debugLevel, &format, args...)
}

// DebugDepth 输出日志
func (lg *Logger) DebugDepth(depth int, args ...any) {
	if lg.disableDebug {
		return
	}
	lg.print(loggerDepth+depth, &debugLevel, args...)
}

// DebugfDepth 输出日志
func (lg *Logger) DebugfDepth(depth int, format string, args ...any) {
	if lg.disableDebug {
		return
	}
	lg.printf(loggerDepth+depth, &debugLevel, &format, args...)
}

// DebugTrace 输出日志
func (lg *Logger) DebugTrace(traceID string, args ...any) {
	if lg.disableDebug {
		return
	}
	if traceID != "" {
		lg.printTrace(loggerDepth, &traceID, &debugLevel, args...)
	} else {
		lg.print(loggerDepth, &debugLevel, args...)
	}
}

// DebugfTrace 输出日志
func (lg *Logger) DebugfTrace(traceID, format string, args ...any) {
	if lg.disableDebug {
		return
	}
	if traceID != "" {
		lg.printfTrace(loggerDepth, &traceID, &debugLevel, &format, args...)
	} else {
		lg.printf(loggerDepth, &debugLevel, &format, args...)
	}
}

// DebugDepthTrace 输出日志
func (lg *Logger) DebugDepthTrace(depth int, traceID string, args ...any) {
	if lg.disableDebug {
		return
	}
	if traceID != "" {
		lg.printTrace(loggerDepth+depth, &traceID, &debugLevel, args...)
	} else {
		lg.print(loggerDepth+depth, &debugLevel, args...)
	}
}

// DebugfDepthTrace 输出日志
func (lg *Logger) DebugfDepthTrace(depth int, traceID, format string, args ...any) {
	if lg.disableDebug {
		return
	}
	if traceID != "" {
		lg.printfTrace(loggerDepth+depth, &traceID, &debugLevel, &format, args...)
	} else {
		lg.printf(loggerDepth+depth, &debugLevel, &format, args...)
	}
}

// IsInfo 返回是否启用 info
func (lg *Logger) IsInfo() bool {
	return !lg.disableInfo
}

// EnableInfo 设置是否启用 info
func (lg *Logger) EnableInfo(enable bool) {
	lg.disableInfo = !enable
}

// Info 输出日志
func (lg *Logger) Info(args ...any) {
	if lg.disableInfo {
		return
	}
	lg.print(loggerDepth, &infoLevel, args...)
}

// Infof 输出日志
func (lg *Logger) Infof(format string, args ...any) {
	if lg.disableInfo {
		return
	}
	lg.printf(loggerDepth, &infoLevel, &format, args...)
}

// InfoDepth 输出日志
func (lg *Logger) InfoDepth(depth int, args ...any) {
	if lg.disableInfo {
		return
	}
	lg.print(loggerDepth+depth, &infoLevel, args...)
}

// InfofDepth 输出日志
func (lg *Logger) InfofDepth(depth int, format string, args ...any) {
	if lg.disableInfo {
		return
	}
	lg.printf(loggerDepth+depth, &infoLevel, &format, args...)
}

// InfoTrace 输出日志
func (lg *Logger) InfoTrace(traceID string, args ...any) {
	if lg.disableInfo {
		return
	}
	if traceID != "" {
		lg.printTrace(loggerDepth, &traceID, &infoLevel, args...)
	} else {
		lg.print(loggerDepth, &infoLevel, args...)
	}
}

// InfofTrace 输出日志
func (lg *Logger) InfofTrace(traceID, format string, args ...any) {
	if lg.disableInfo {
		return
	}
	if traceID != "" {
		lg.printfTrace(loggerDepth, &traceID, &infoLevel, &format, args...)
	} else {
		lg.printf(loggerDepth, &infoLevel, &format, args...)
	}
}

// InfoDepthTrace 输出日志
func (lg *Logger) InfoDepthTrace(depth int, traceID string, args ...any) {
	if lg.disableInfo {
		return
	}
	if traceID != "" {
		lg.printTrace(loggerDepth+depth, &traceID, &infoLevel, args...)
	} else {
		lg.print(loggerDepth+depth, &infoLevel, args...)
	}
}

// InfofDepthTrace 输出日志
func (lg *Logger) InfofDepthTrace(depth int, traceID, format string, args ...any) {
	if lg.disableInfo {
		return
	}
	if traceID != "" {
		lg.printfTrace(loggerDepth+depth, &traceID, &infoLevel, &format, args...)
	} else {
		lg.printf(loggerDepth+depth, &infoLevel, &format, args...)
	}
}

// IsWarn 返回是否启用 warn
func (lg *Logger) IsWarn() bool {
	return !lg.disableWarn
}

// EnableWarn 设置是否启用 warn
func (lg *Logger) EnableWarn(enable bool) {
	lg.disableWarn = !enable
}

// Warn 输出日志
func (lg *Logger) Warn(args ...any) {
	if lg.disableWarn {
		return
	}
	lg.print(loggerDepth, &warnLevel, args...)
}

// Warnf 输出日志
func (lg *Logger) Warnf(format string, args ...any) {
	if lg.disableWarn {
		return
	}
	lg.printf(loggerDepth, &warnLevel, &format, args...)
}

// WarnDepth 输出日志
func (lg *Logger) WarnDepth(depth int, args ...any) {
	if lg.disableWarn {
		return
	}
	lg.print(loggerDepth+depth, &warnLevel, args...)
}

// WarnfDepth 输出日志
func (lg *Logger) WarnfDepth(depth int, format string, args ...any) {
	if lg.disableWarn {
		return
	}
	lg.printf(loggerDepth+depth, &warnLevel, &format, args...)
}

// WarnTrace 输出日志
func (lg *Logger) WarnTrace(traceID string, args ...any) {
	if lg.disableWarn {
		return
	}
	if traceID != "" {
		lg.printTrace(loggerDepth, &traceID, &warnLevel, args...)
	} else {
		lg.print(loggerDepth, &warnLevel, args...)
	}
}

// WarnfTrace 输出日志
func (lg *Logger) WarnfTrace(traceID, format string, args ...any) {
	if lg.disableWarn {
		return
	}
	if traceID != "" {
		lg.printfTrace(loggerDepth, &traceID, &warnLevel, &format, args...)
	} else {
		lg.printf(loggerDepth, &warnLevel, &format, args...)
	}
}

// WarnDepthTrace 输出日志
func (lg *Logger) WarnDepthTrace(depth int, traceID string, args ...any) {
	if lg.disableWarn {
		return
	}
	if traceID != "" {
		lg.printTrace(loggerDepth+depth, &traceID, &warnLevel, args...)
	} else {
		lg.print(loggerDepth+depth, &warnLevel, args...)
	}
}

// WarnfDepthTrace 输出日志
func (lg *Logger) WarnfDepthTrace(depth int, traceID, format string, args ...any) {
	if lg.disableWarn {
		return
	}
	if traceID != "" {
		lg.printfTrace(loggerDepth+depth, &traceID, &warnLevel, &format, args...)
	} else {
		lg.printf(loggerDepth+depth, &warnLevel, &format, args...)
	}
}

// IsError 返回是否启用 error
func (lg *Logger) IsError() bool {
	return !lg.disableError
}

// EnableError 设置是否启用 error
func (lg *Logger) EnableError(enable bool) {
	lg.disableError = !enable
}

// Error 输出日志
func (lg *Logger) Error(args ...any) {
	if lg.disableError {
		return
	}
	lg.print(loggerDepth, &errorLevel, args...)
}

// Errorf 输出日志
func (lg *Logger) Errorf(format string, args ...any) {
	if lg.disableError {
		return
	}
	lg.printf(loggerDepth, &errorLevel, &format, args...)
}

// ErrorDepth 输出日志
func (lg *Logger) ErrorDepth(depth int, args ...any) {
	if lg.disableError {
		return
	}
	lg.print(loggerDepth+depth, &errorLevel, args...)
}

// ErrorfDepth 输出日志
func (lg *Logger) ErrorfDepth(depth int, format string, args ...any) {
	if lg.disableError {
		return
	}
	lg.printf(loggerDepth+depth, &errorLevel, &format, args...)
}

// ErrorTrace 输出日志
func (lg *Logger) ErrorTrace(traceID string, args ...any) {
	if lg.disableError {
		return
	}
	if traceID != "" {
		lg.printTrace(loggerDepth, &traceID, &errorLevel, args...)
	} else {
		lg.print(loggerDepth, &errorLevel, args...)
	}
}

// ErrorfTrace 输出日志
func (lg *Logger) ErrorfTrace(traceID, format string, args ...any) {
	if lg.disableError {
		return
	}
	if traceID != "" {
		lg.printfTrace(loggerDepth, &traceID, &errorLevel, &format, args...)
	} else {
		lg.printf(loggerDepth, &errorLevel, &format, args...)
	}
}

// ErrorDepthTrace 输出日志
func (lg *Logger) ErrorDepthTrace(depth int, traceID string, args ...any) {
	if lg.disableError {
		return
	}
	if traceID != "" {
		lg.printTrace(loggerDepth+depth, &traceID, &errorLevel, args...)
	} else {
		lg.print(loggerDepth+depth, &errorLevel, args...)
	}
}

// ErrorfDepthTrace 输出日志
func (lg *Logger) ErrorfDepthTrace(depth int, traceID, format string, args ...any) {
	if lg.disableError {
		return
	}
	if traceID != "" {
		lg.printfTrace(loggerDepth+depth, &traceID, &errorLevel, &format, args...)
	} else {
		lg.printf(loggerDepth+depth, &errorLevel, &format, args...)
	}
}
