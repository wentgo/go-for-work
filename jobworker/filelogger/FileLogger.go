package filelogger

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	EMERG   = iota // 0) system is unusable
	ALERT          // 1) action must be taken immediately
	CRIT           // 2) critical conditions
	ERROR          // 3) error conditions
	WARNING        // 4) warning conditions
	NOTICE         // 5) normal, but significant, condition
	INFO           // 6) informational message
	DEBUG          // 7) debug-level message
)

type FileLogger struct {
	filename string
	buffer   *bytes.Buffer
	logger   *log.Logger
	level    int
}

func New(fname string) *FileLogger {
	buf := &bytes.Buffer{}
	logger := log.New(buf, "", 0)
	return &FileLogger{
		filename: fname,
		buffer:   buf,
		logger:   logger,
		level:    DEBUG,
	}
}

func (this *FileLogger) Disable() {
	this.logger.SetOutput(ioutil.Discard)
}

func (this *FileLogger) Enable() {
	this.logger.SetOutput(this.buffer)
}

func (this *FileLogger) SetLevel(level int) {
	this.level = level
}

// internal use only
func (this *FileLogger) write(level int, format string, v ...interface{}) {
	if level > this.level {
		return
	}

	text := this.getLevelText(level)
	now := time.Now().Format("2006-01-02 15:04:05")

	var fmtstr string
	if text != "" {
		fmtstr = fmt.Sprintf("%s [%s] %s", now, text, format)
	} else {
		fmtstr = fmt.Sprintf("%s %s", now, format)
	}

	this.logger.Printf(fmtstr, v...)
}

// internal use only
func (this *FileLogger) getLevelText(level int) string {
	text := ""
	switch level {
	case EMERG:
		text = "EMERG"
	case ALERT:
		text = "ALERT"
	case CRIT:
		text = "CRIT"
	case ERROR:
		text = "ERROR"
	case WARNING:
		text = "WARNING"
	case NOTICE:
		text = "NOTICE"
	case INFO:
		text = "INFO"
	case DEBUG:
		text = "DEBUG"
	}
	return text
}

func (this *FileLogger) Print(format string, v ...interface{}) {
	this.write(-1, format, v...)
}

func (this *FileLogger) Emerg(format string, v ...interface{}) {
	this.write(EMERG, format, v...)
}

func (this *FileLogger) Alert(format string, v ...interface{}) {
	this.write(ALERT, format, v...)
}

func (this *FileLogger) Crit(format string, v ...interface{}) {
	this.write(CRIT, format, v...)
}

func (this *FileLogger) Error(format string, v ...interface{}) {
	this.write(ERROR, format, v...)
}

func (this *FileLogger) Warning(format string, v ...interface{}) {
	this.write(WARNING, format, v...)
}

func (this *FileLogger) Notice(format string, v ...interface{}) {
	this.write(NOTICE, format, v...)
}

func (this *FileLogger) Info(format string, v ...interface{}) {
	this.write(INFO, format, v...)
}

func (this *FileLogger) Debug(format string, v ...interface{}) {
	this.write(DEBUG, format, v...)
}

func (this *FileLogger) GetFilename() string {
	ext := filepath.Ext(this.filename)
	name := strings.TrimSuffix(this.filename, ext)
	return name + "-" + time.Now().Format("2006-01-02") + ext
}

func (this *FileLogger) Flush() {
	if this.buffer.Len() == 0 {
		return
	}

	filename := this.GetFilename()
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Fprint(file, this.buffer)
	this.buffer.Reset()
}
