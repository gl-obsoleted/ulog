package core

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type LogWriter struct {
	LogFile *os.File
}

func (lw LogWriter) Write(p []byte) (int, error) {
	if _, err := fmt.Print(string(p)); err != nil {
		return 0, err
	}

	if _, err := lw.LogFile.Write(p); err != nil {
		fmt.Print(err.Error())
		return 0, err
	}

	return len(p), nil
}

var GLogWriter LogWriter

var GlobalErrorCount int = 0
var GlobalFatalErrorCount int = 0

var GLogDebug = flag.Bool("d", false, "turn on debug output")

func InitLogging() error {
	t := time.Now()
	logging_dir := fmt.Sprintf("temp/%02d-%02d-%02d/", t.Year(), t.Month(), t.Day())
	if err := CreateDirIfNotExists(logging_dir); err != nil {
		return NewStdErr(ERR_EnvFatal, err.Error())
	}
	log.Println("logging directory initialized: " + logging_dir)

	logging_filename := fmt.Sprintf("%02d-%02d-%02d-%02d-%02d-%02d.log", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	dst, err := os.Create(logging_dir + logging_filename)
	if err != nil {
		return NewStdErr(ERR_EnvFatal, err.Error())
	}

	// GLogWriter is ready to go, assemble into go runtime logging
	GLogWriter = LogWriter{}
	GLogWriter.LogFile = dst

	log.SetFlags(log.Ldate | log.Ldate | log.Lshortfile)
	log.SetOutput(GLogWriter)
	log.Println("--- logging system initialized ---")
	log.Println("  Log File: ", GLogWriter.LogFile.Name())
	log.Println("  Debug Output: ", *GLogDebug)

	return nil
}

func DestroyLogging() {
	log.Println("Destroying logging...")
	GLogWriter.LogFile.Close()
}

// =========== Public Functions ===========

func LogDebug(format string, args ...interface{}) {
	if !(*GLogDebug) {
		return
	}

	log.Printf("{debug} "+format, args...)
}

func LogErrorf(format string, args ...interface{}) {
	GlobalErrorCount++
	log.Printf("{error(#"+strconv.Itoa(GlobalErrorCount)+")} "+format, args...)
}

func LogErrorln(args ...interface{}) {
	GlobalErrorCount++
	log.Printf("{error(#%v)} %s\n", GlobalErrorCount, fmt.Sprint(args...))
}

func LogError(title string, err error) {
	if err == nil {
		log.Println("\nLogError() is requested but 'nil' error is passed in, ignored.")
		return
	}

	GlobalErrorCount++
	log.Printf("{error(#%v)} %s - (%s)\n", GlobalErrorCount, title, err.Error())
}

func LogFatalError(title string, err error) {
	if err == nil {
		log.Println("\nLogFatalError() is requested but 'nil' error is passed in, ignored.")
		return
	}

	GlobalFatalErrorCount++
	log.Printf("{fatal(#%v)} %s - (%s)\n", GlobalFatalErrorCount, title, err.Error())
}
