package helpers

import (
	"bytes"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/iakrevetkho/components-tests/cott/config"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(config *config.Config) error {
	SetLoggerFormat(config)

	// Check that app can write the log file
	logFile, err := os.Create(config.Log.FilePath)
	if err != nil {
		return err
	}
	logFile.Close()

	// Create file for rotated logs
	rotatedLog := &lumberjack.Logger{
		Filename:   config.Log.FilePath,
		MaxSize:    config.Log.MaxFileSizeInMb,
		MaxAge:     config.Log.MaxFileAgeInDays,
		MaxBackups: config.Log.MaxFilesCount,
		Compress:   config.Log.CompressOldFiles,
	}

	// Create daily task with Cron to rotate logs by timestamp
	c := cron.New()
	if err := c.AddFunc("@daily", func() {
		rotateLogFileIfNotEmpty(config.Log.FilePath, rotatedLog)
	}); err != nil {
		return err
	}
	c.Start()

	// Create writer into the console and file simultaneously
	mw := io.MultiWriter(os.Stdout, rotatedLog)
	logrus.SetOutput(mw)

	return nil
}

func SetLoggerFormat(config *config.Config) {
	// Set logger formatter
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			functionName := path.Base(f.Function)

			var funcNameBuf bytes.Buffer
			funcNameBuf.WriteString(functionName)
			funcNameBuf.WriteString("()")

			var filePathBuf bytes.Buffer
			filePathBuf.WriteByte('\t')
			filePathBuf.WriteString(filename)
			filePathBuf.WriteByte(':')
			filePathBuf.WriteString(strconv.FormatInt(int64(f.Line), 10))

			return funcNameBuf.String(), filePathBuf.String()
		},
	})
	logrus.SetLevel(config.Log.Level)
}

func rotateLogFileIfNotEmpty(logFilePath string, rotatedLog *lumberjack.Logger) {
	// Get file's stats to check that log file is not empty
	rotatedLogFileStat, err := os.Stat(logFilePath)
	if err != nil {
		logrus.WithError(err).Error("Couldn't read log file stats")
		return
	}
	// Check that log file is not empty
	rotatedLogFileSize := rotatedLogFileStat.Size()
	if rotatedLogFileSize > 0 {
		// Rotate non empty log file
		logrus.WithField("fileSize", rotatedLogFileSize).Debug("Rotate log file")
		if err := rotatedLog.Rotate(); err != nil {
			logrus.WithError(err).Error("Couldn't rotate log file")
			return
		}
	}
}
