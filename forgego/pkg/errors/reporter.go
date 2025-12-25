package errors

import (
	"fmt"
	"log"
	"os"
)

type LogReporter struct {
	logger *log.Logger
}

func NewLogReporter() *LogReporter {
	return &LogReporter{
		logger: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
	}
}

func (lr *LogReporter) Report(err *Error) error {
	lr.logger.Printf("Code: %s, Message: %s, Status: %d", err.Code, err.Message, err.Status)
	
	if err.Err != nil {
		lr.logger.Printf("Caused by: %v", err.Err)
	}
	
	if len(err.Context) > 0 {
		lr.logger.Printf("Context: %+v", err.Context)
	}
	
	if len(err.Stack) > 0 {
		lr.logger.Printf("Stack trace:")
		for _, frame := range err.Stack {
			lr.logger.Printf("  %s", frame)
		}
	}
	
	return nil
}

type FileReporter struct {
	file   *os.File
	logger *log.Logger
}

func NewFileReporter(filePath string) (*FileReporter, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	
	return &FileReporter{
		file:   file,
		logger: log.New(file, "[ERROR] ", log.LstdFlags),
	}, nil
}

func (fr *FileReporter) Report(err *Error) error {
	fr.logger.Printf("Code: %s, Message: %s, Status: %d", err.Code, err.Message, err.Status)
	
	if err.Err != nil {
		fr.logger.Printf("Caused by: %v", err.Err)
	}
	
	if len(err.Context) > 0 {
		fr.logger.Printf("Context: %+v", err.Context)
	}
	
	if len(err.Stack) > 0 {
		fr.logger.Printf("Stack trace:")
		for _, frame := range err.Stack {
			fr.logger.Printf("  %s", frame)
		}
	}
	
	return nil
}

func (fr *FileReporter) Close() error {
	if fr.file != nil {
		return fr.file.Close()
	}
	return nil
}

type JSONReporter struct {
	encoder func(*Error) string
}

func NewJSONReporter() *JSONReporter {
	return &JSONReporter{
		encoder: func(err *Error) string {
			return fmt.Sprintf(`{"code":"%s","message":"%s","status":%d,"timestamp":"%s"}`, 
				err.Code, err.Message, err.Status, err.Timestamp)
		},
	}
}

func (jr *JSONReporter) Report(err *Error) error {
	fmt.Println(jr.encoder(err))
	return nil
}

