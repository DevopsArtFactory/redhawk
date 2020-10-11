package tools

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/DevopsArtFactory/redhawk/pkg/constants"
)

// ClearOsEnv removes all environment variables about AWS
func ClearOsEnv() error {
	logrus.Debugf("remove environment variable")
	if err := os.Unsetenv("AWS_ACCESS_KEY_ID"); err != nil {
		return err
	}
	if err := os.Unsetenv("AWS_SECRET_ACCESS_KEY"); err != nil {
		return err
	}

	if err := os.Unsetenv("AWS_SESSION_TOKEN"); err != nil {
		return err
	}

	return nil
}

// Check if file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

//Figure out if string is in array
func IsStringInArray(s string, arr []string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}

	return false
}

// IsExpired compares current time with (targetDate + timeAdded)
func IsExpired(targetDate time.Time, timeAdded time.Duration) bool {
	return time.Since(targetDate.Add(timeAdded)) > 0
}

// SetUpLogs set logrus log format
func SetUpLogs(stdErr io.Writer, level string) error {
	logrus.SetOutput(stdErr)
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("parsing log level: %w", err)
	}
	logrus.SetLevel(lvl)
	return nil
}

// CheckValidFormat checks if format is valid or not
func CheckValidFormat(format string) error {
	if !IsStringInArray(format, constants.ValidFormats) {
		return fmt.Errorf("format is not supported: %s", format)
	}

	return nil
}
