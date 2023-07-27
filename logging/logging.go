package logging

import (
	"fmt"
	"os"
	"time"

	"dynamicledger.com/testnet-deployer/structs"
	alsdk "github.com/activeledger/SDK-Golang/v2"
)

type Logger struct {
	logFolder      string
	writeToFile    bool
	verboseLogging bool
	headlessMode   bool
	logFileName    string
}

type level int

const (
	ERR   level = 0
	FATAL level = 1
	INFO  level = 2
	WARN  level = 3
	DEBUG level = 4
)

func CreateLogger() Logger {
	return Logger{
		logFolder:      "",
		writeToFile:    false,
		verboseLogging: true,
		headlessMode:   true,
		logFileName:    "",
	}
}

func (l *Logger) SetConfig(config *structs.Config) {
	l.logFolder = config.LogFolder
	l.writeToFile = config.LogToFile
	l.verboseLogging = config.VerboseLogging
	l.headlessMode = config.HeadlessMode
}

func (l *Logger) GetUserInput(msg string) string {

	if l.headlessMode {
		var blank string
		return blank
	}

	fmt.Print(msg)

	var resp string
	fmt.Scanln(&resp)

	return resp
}

func (l *Logger) Info(msg string) {
	l.print(msg, INFO, nil)
}

func (l *Logger) Warn(msg string) {
	l.print(msg, WARN, nil)
}

func (l *Logger) Debug(msg string) {
	l.print(msg, DEBUG, nil)
}

func (l *Logger) Error(err error, msg string) {
	l.print(msg, ERR, err)
}

func (l *Logger) Fatal(err error, msg string) {
	l.print(msg, FATAL, err)
	os.Exit(1)
}

func (l *Logger) ActiveledgerError(err error, resp alsdk.Response, msg string) {
	l.handleALError(err, resp, msg)
	os.Exit(2)
}

func (l *Logger) print(msg string, lv level, err error) {
	output := "\u001b[37m;1m"

	timestamp := getTimestamp()

	switch lv {
	case ERR:
		output = fmt.Sprintf(
			"%s\u001b[41m[%s - Error]\u001b[0m: %s\n%s\n",
			output,
			timestamp,
			msg,
			err,
		)

	case FATAL:
		output = fmt.Sprintf(
			"%s\u001b[41m[%s - FATAL Error]\u001b[0m: %s\n%s\n\nFatal error, shutting down",
			output,
			timestamp,
			msg,
			err,
		)

	case INFO:
		output = fmt.Sprintf(
			"%s\u001b[42m[%s - Info]\u001b[0m: %s\n",
			output,
			timestamp,
			msg,
		)

	case WARN:
		output = fmt.Sprintf(
			"%s\u001b[44m[%s - Warning]\u001b[0m: %s\n",
			output,
			timestamp,
			msg,
		)

	case DEBUG:
		output = fmt.Sprintf(
			"%s\u001b[46m[%s - Debug]\u001b[0m: %s\n",
			output,
			timestamp,
			msg,
		)
	}

	if !l.headlessMode {
		fmt.Print(output)
	}

	// If verbose logging off skip outputting debugs
	if !l.verboseLogging && lv == DEBUG {
		return
	}

	if l.writeToFile {
		l.writeFile(output)
	}
}

func (l *Logger) writeFile(data string) {
	path := l.logFolder
	timestamp := getTimestamp()

	if l.logFileName == "" {
		l.logFileName = fmt.Sprintf("%s.log", timestamp)
	}

	filepath := fmt.Sprintf("%s/%s", path, l.logFileName)

	f, err := os.OpenFile(
		filepath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		output := fmt.Sprintf(
			"\u001b[37m;1m\u001b[41m[%s - Error]\u001b[0m: Unable to open log file \"%s\"\n%s\n\n",
			timestamp,
			filepath,
			err,
		)

		fmt.Print(output)
	}

	defer f.Close()

	if _, err := f.WriteString(data); err != nil {
		output := fmt.Sprintf(
			"\u001b[37m;1m\u001b[41m[%s - Error]\u001b[0m: Unable to write to log file \"%s\"\n%s\n\n",
			timestamp,
			filepath,
			err,
		)

		fmt.Print(output)
	}
}

func getTimestamp() string {
	timestamp := time.Now()

	return timestamp.Format(time.RFC3339)
}
