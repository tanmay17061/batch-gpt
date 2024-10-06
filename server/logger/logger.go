package logger

import (
    "log"
    "os"
)

var (
    InfoLogger  *log.Logger
    WarnLogger  *log.Logger
    ErrorLogger *log.Logger
)

func init() {
    InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    WarnLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
    ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}