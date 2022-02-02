package domain

import (
	"github.com/sirupsen/logrus"
)

type Config struct {
	Log       LogConfig
	Report    ReportConfig
	TestCases []TestCase
}

type LogConfig struct {
	Level            logrus.Level `default:"info" env:"LOG_LEVEL"`
	FilePath         string       `default:"/var/log/cott/cott.log" env:"LOG_FILE_PATH"`
	MaxFileSizeInMb  int          `default:"10" env:"LOG_MAX_FILE_SIZE_IN_MB"`
	MaxFilesCount    int          `default:"7" env:"LOG_MAX_FILES_COUNT"`
	MaxFileAgeInDays int          `default:"7" env:"LOG_MAX_FILE_AGE_IN_DAYS"`
	CompressOldFiles bool         `default:"true" env:"LOG_COMPRESS_OLD_FILES"`
}

type ReportConfig struct {
	FilePath string `default:"report.json" env:"REPORT_FILE_PATH"`
}
