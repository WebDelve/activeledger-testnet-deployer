package logging

import (
	"errors"
	"fmt"
	"os"
	"strings"
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

// Longest prefix, used for padding
const prefixLength = 47

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

	fmt.Printf("\x1b[36;1m[Input] %s\x1b[0m", msg)

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

func addPadding(s string) string {
	strLen := len(s)
	padding := prefixLength - strLen

	spl := strings.Split(s, " - ")

	if padding == 0 {
		s = fmt.Sprintf("%s %s] \x1b[0m", spl[0], spl[1])

	} else {
		for i := 0; i <= padding; i++ {
			spl[0] = spl[0] + " "
		}

		s = spl[0] + spl[1] + "] \x1b[0m"
	}

	return s
}

func (l *Logger) print(msg string, lv level, err error) {
	output := ""
	errMsg := ""

	timestamp := getTimestamp()

	switch lv {
	case ERR:
		prefix := fmt.Sprintf("\x1b[31;1m[%s - Error", timestamp)
		prefix = addPadding(prefix)

		output = fmt.Sprintf(
			"%s%s\n",
			prefix,
			msg,
		)
		errMsg = fmt.Sprintf(
			"%s%s\n",
			prefix,
			err,
		)

	case FATAL:
		prefix := fmt.Sprintf("\x1b[31;1m[%s - FATAL Error", timestamp)
		prefix = addPadding(prefix)

		output = fmt.Sprintf(
			"%s%s\n",
			prefix,
			msg,
		)
		errMsg = fmt.Sprintf(
			"%s%s\n",
			prefix,
			err,
		)

		errMsg = errMsg + "\n\nFatal error, shutting down\n"

	case INFO:
		prefix := fmt.Sprintf("\x1b[32;1m[%s - Info", timestamp)
		prefix = addPadding(prefix)

		output = fmt.Sprintf(
			"%s%s\n",
			prefix,
			msg,
		)

	case WARN:
		prefix := fmt.Sprintf("\x1b[33;1m[%s - Warning", timestamp)
		prefix = addPadding(prefix)

		output = fmt.Sprintf(
			"%s%s\n",
			prefix,
			msg,
		)

	case DEBUG:
		prefix := fmt.Sprintf("\x1b[34;1m[%s - Debug", timestamp)
		prefix = addPadding(prefix)

		output = fmt.Sprintf(
			"%s%s\n",
			prefix,
			msg,
		)
	}

	// If verbose logging off skip outputting debugs
	if !l.verboseLogging && lv == DEBUG {
		return
	}

	if !l.headlessMode {
		fmt.Print(output)

		if lv == ERR {
			fmt.Print(errMsg)
		}
	}

	if l.writeToFile {
		data := cleanLogForFile(output, lv, timestamp)
		l.writeFile(data)
	}
}

func cleanLogForFile(data string, lv level, ts string) string {
	split := strings.Split(data, "] ")

	levelString := ""
	switch lv {
	case FATAL:
		levelString = "FATAL Error"

	case ERR:
		levelString = "Error"

	case INFO:
		levelString = "Info"

	case WARN:
		levelString = "Warn"

	case DEBUG:
		levelString = "Debug"

	}

	output := fmt.Sprintf("[%s %s] %s", ts, levelString, split[1])
	return output
}

func (l *Logger) writeFile(data string) {
	path := l.logFolder
	timestamp := getTimestamp()

	if l.logFileName == "" {
		l.logFileName = fmt.Sprintf("%s.log", timestamp)
	}

	if err := folderCheck(l.logFolder); err != nil {
		prefix := fmt.Sprintf("\x1b[31;1m[%s - FATAL Error", timestamp)
		prefix = addPadding(prefix)

		msg := "Error checking for log folder"

		fmt.Printf(
			"%s%s\n",
			prefix,
			msg,
		)
		fmt.Printf(
			"%s%s\n\nFatal error, shutting down\n",
			prefix,
			err,
		)

	}

	filepath := fmt.Sprintf("%s/%s", path, l.logFileName)

	f, err := os.OpenFile(
		filepath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		prefix := fmt.Sprintf("\x1b[31;1m[%s - FATAL Error", timestamp)
		prefix = addPadding(prefix)

		msg := "Unable to open log file"

		fmt.Printf(
			"%s%s\n",
			prefix,
			msg,
		)
		fmt.Printf(
			"%s%s\n\nFatal error, shutting down\n",
			prefix,
			err,
		)

		os.Exit(3)
	}

	defer f.Close()

	if _, err := f.WriteString(data); err != nil {
		prefix := fmt.Sprintf("\x1b[31;1m[%s - FATAL Error", timestamp)
		prefix = addPadding(prefix)

		msg := "Unable to write to log file"

		fmt.Printf(
			"%s%s\n",
			prefix,
			msg,
		)
		fmt.Printf(
			"%s%s\n\nFatal error, shutting down\n",
			prefix,
			err,
		)

		os.Exit(3)
	}
}

func getTimestamp() string {
	timestamp := time.Now()

	return timestamp.Format(time.RFC3339)
}

func folderCheck(folder string) error {
	_, err := os.Stat(folder)

	if errors.Is(err, os.ErrNotExist) {
		if err = os.Mkdir(folder, 0755); err != nil {
			return err
		}
	}

	return err
}
