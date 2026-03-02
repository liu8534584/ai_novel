package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	LLMLogger   *log.Logger
	logFile     *os.File
)

func Init(level, filename string) error {
	flags := log.Ldate | log.Ltime | log.Lshortfile

	// 确保日志目录存在
	dir := filepath.Dir(filename)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}
	}

	// 打开日志文件
	var err error
	logFile, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// 设置多路输出
	infoWriter := io.MultiWriter(os.Stdout, logFile)
	errorWriter := io.MultiWriter(os.Stderr, logFile)
	llmWriter := io.MultiWriter(os.Stdout, logFile)

	InfoLogger = log.New(infoWriter, "[INFO] ", flags)
	ErrorLogger = log.New(errorWriter, "[ERROR] ", flags)
	LLMLogger = log.New(llmWriter, "[LLM] ", flags)

	return nil
}

func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

func init() {
	// 默认初始化，防止在 Init 调用前被使用
	flags := log.Ldate | log.Ltime | log.Lshortfile
	InfoLogger = log.New(os.Stdout, "[INFO] ", flags)
	ErrorLogger = log.New(os.Stderr, "[ERROR] ", flags)
	LLMLogger = log.New(os.Stdout, "[LLM] ", flags)
}

func Info(format string, v ...interface{}) {
	InfoLogger.Output(2, fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	ErrorLogger.Output(2, fmt.Sprintf(format, v...))
}

func LLM(format string, v ...interface{}) {
	LLMLogger.Output(2, fmt.Sprintf(format, v...))
}

func LogLLMRequest(bookID uint, model string, messages interface{}, response interface{}, err error) {
	status := "SUCCESS"
	if err != nil {
		status = "FAILED"
	}

	LLM("BookID: %d | Model: %s | Status: %s\n--- [REQUEST] ---\n%+v\n--- [RESPONSE] ---\n%+v\n--- [ERROR] ---\n%v\n----------------",
		bookID, model, status, messages, response, err)
}
