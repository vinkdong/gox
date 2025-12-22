package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

const ClrEnd = "\x1b[0m"
const ClrDEBUG = "\x1b[37;2m"
const ClrSuccess = "\x1b[32;2m"
const ClrWarn = "\x1b[33;2m"

var (
	debug        = false
	hideFileInfo = true

	std           *Logger
	fileWriter    io.Writer // 纯文本文件 writer（无颜色）
	consoleWriter io.Writer // 带颜色的控制台 writer
	mu            sync.Mutex

	fileEnabled    bool
	consoleEnabled = true
)

type Logger struct {
	mu     sync.Mutex
	prefix string
	flag   int
	out    io.Writer // 最终输出目标（MultiWriter）
	buf    []byte
}

func init() {
	// 初始只输出到控制台（带颜色）
	consoleWriter = os.Stdout
	std = New(consoleWriter, "", log.Ldate|log.Ltime)
}

// ==================== 配置函数 ====================
func SetFileOutput(options ...FileOption) {
	mu.Lock()
	defer mu.Unlock()

	// 初始化 lumberjack
	lj := &lumberjack.Logger{
		Filename:   "/var/log/gox.log",
		MaxSize:    100, // MB
		MaxBackups: 10,
		MaxAge:     180, // days
		Compress:   true,
	}

	for _, opt := range options {
		opt(lj)
	}

	// 确保目录存在
	dir := filepath.Dir(lj.Filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		// 如果目录都创建不了，至少打印警告（但不会崩溃）
		fmt.Fprintf(os.Stderr, "log: 创建日志目录失败 %s: %v\n", dir, err)
	}

	fileWriter = lj
	fileEnabled = true
	consoleEnabled = false // 严格满足第3点：启用文件后关闭控制台

	updateOutputWriter()
}

// 更新最终输出目标
func updateOutputWriter() {
	var writers []io.Writer
	if fileEnabled && fileWriter != nil {
		writers = append(writers, fileWriter) // 文件：纯文本
	}
	if consoleEnabled && consoleWriter != nil {
		writers = append(writers, consoleWriter) // 控制台：带颜色
	}

	if len(writers) == 0 {
		std.out = io.Discard
	} else if len(writers) == 1 {
		std.out = writers[0]
	} else {
		std.out = io.MultiWriter(writers...)
	}
}

// FileOption 配置项
type FileOption func(*lumberjack.Logger)

func WithFilename(filename string) FileOption {
	return func(l *lumberjack.Logger) { l.Filename = filename }
}
func WithMaxSize(mb int) FileOption {
	return func(l *lumberjack.Logger) { l.MaxSize = mb }
}
func WithMaxBackups(num int) FileOption {
	return func(l *lumberjack.Logger) { l.MaxBackups = num }
}
func WithMaxAge(days int) FileOption {
	return func(l *lumberjack.Logger) { l.MaxAge = days }
}
func WithCompress(compress bool) FileOption {
	return func(l *lumberjack.Logger) { l.Compress = compress }
}

func EnableConsole() {
	mu.Lock()
	defer mu.Unlock()
	consoleEnabled = true
	updateOutputWriter()
}

func DisableConsole() {
	mu.Lock()
	defer mu.Unlock()
	consoleEnabled = false
	updateOutputWriter()
}

// EnableWrite
func EnableWrite() { /* 已废弃，保留兼容 */ }
func SetFilename(name string) {
	SetFileOutput(WithFilename(name))
}

func EnableDebug()    { debug = true }
func EnableFileInfo() { hideFileInfo = false }

// ==================== Logger 核心 ====================

func New(out io.Writer, prefix string, flag int) *Logger {
	return &Logger{out: out, prefix: prefix, flag: flag}
}

func (l *Logger) Output(calldepth int, s string, color string) error {
	now := time.Now()
	var file string
	var line int

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.flag&(log.Lshortfile|log.Llongfile) != 0 {
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}

	l.buf = l.buf[:0]

	// 关键：只有控制台输出才加颜色，文件输出不加
	if consoleEnabled && !fileEnabled { // 纯控制台模式
		if color != "" {
			l.buf = append(l.buf, color...)
		}
	}

	l.formatHeader(&l.buf, now, file, line)
	l.buf = append(l.buf, s...)

	if consoleEnabled && !fileEnabled {
		if color != "" {
			l.buf = append(l.buf, ClrEnd...)
		}
	}

	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}

	_, err := l.out.Write(l.buf)
	return err
}

// 精简版工具函数
func itoa(buf *[]byte, i int, wid int) {
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func (l *Logger) formatHeader(buf *[]byte, t time.Time, file string, line int) {
	*buf = append(*buf, l.prefix...)
	if l.flag&(log.Ldate|log.Ltime) != 0 {
		year, month, day := t.Date()
		itoa(buf, year, 4)
		*buf = append(*buf, '/')
		itoa(buf, int(month), 2)
		*buf = append(*buf, '/')
		itoa(buf, day, 2)
		*buf = append(*buf, ' ')
		hour, min, sec := t.Clock()
		itoa(buf, hour, 2)
		*buf = append(*buf, ':')
		itoa(buf, min, 2)
		*buf = append(*buf, ':')
		itoa(buf, sec, 2)
		*buf = append(*buf, ' ')
	}
	if l.flag&(log.Lshortfile|log.Llongfile) != 0 {
		if l.flag&log.Lshortfile != 0 {
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					file = file[i+1:]
					break
				}
			}
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}

// ==================== 公共日志函数（保持原调用方式） ====================

func getLineInfo(skip bool) string {
	if skip {
		return ""
	}
	_, file, line, _ := runtime.Caller(2)
	parts := strings.Split(file, "/")
	if len(parts) >= 2 {
		return fmt.Sprintf("%s/%s:%d ", parts[len(parts)-2], parts[len(parts)-1], line)
	}
	return fmt.Sprintf("%s:%d ", file, line)
}

// 统一入口宏
func logOutput(s string, color string) {
	std.Output(2, s, color)
}

func Info(v interface{})               { logOutput(fmt.Sprintf("[INFO] %s%v", getLineInfo(hideFileInfo), v), "") }
func Infof(f string, a ...interface{}) { Info(fmt.Sprintf(f, a...)) }

func Success(v interface{}) {
	logOutput(fmt.Sprintf("[INFO] %s%v", getLineInfo(hideFileInfo), v), ClrSuccess)
}
func Successf(f string, a ...interface{}) { Success(fmt.Sprintf(f, a...)) }

func Warn(v interface{}) {
	logOutput(fmt.Sprintf("[WARN] %s%v", getLineInfo(hideFileInfo), v), ClrWarn)
}
func Warnf(f string, a ...interface{}) { Warn(fmt.Sprintf(f, a...)) }

func Error(v interface{})               { logOutput(fmt.Sprintf("[ERROR] %s%v", getLineInfo(hideFileInfo), v), "") }
func Errorf(f string, a ...interface{}) { Error(fmt.Sprintf(f, a...)) }

func ErrorLine(v interface{})               { logOutput(fmt.Sprintf("[ERROR] %s%v", getLineInfo(false), v), "") }
func ErrorLinef(f string, a ...interface{}) { ErrorLine(fmt.Sprintf(f, a...)) }

func Debug(v interface{}) {
	if debug {
		logOutput(fmt.Sprintf("[DEBUG] %s%v", getLineInfo(hideFileInfo), v), ClrDEBUG)
	}
}
func Debugf(f string, a ...interface{}) {
	if debug {
		Debug(fmt.Sprintf(f, a...))
	}
}

// 兼容旧 Write
func Write(v interface{})               { Info(v) }
func Writef(f string, a ...interface{}) { Write(fmt.Sprintf(f, a...)) }

func Lock()   { std.mu.Lock() }
func Unlock() { std.mu.Unlock() }
