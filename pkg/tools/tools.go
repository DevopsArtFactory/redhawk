/*
Copyright 2020 The redhawk Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tools

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"reflect"
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

// DecodeURLEncodedString decodes url-encoded string
func DecodeURLEncodedString(encoded string) (string, error) {
	decoded, err := url.QueryUnescape(encoded)
	if err != nil {
		return constants.EmptyString, err
	}

	return decoded, nil
}

// Formatting removes nil value
func Formatting(i interface{}) interface{} {
	switch reflect.TypeOf(i).String() {
	case "*time.Time":
		if i.(*time.Time) == nil {
			return "-"
		}
		return i.(*time.Time).Format("2006-01-02 15:04:05")
	case "*string":
		if i.(*string) == nil {
			return "-"
		}
	case "*int":
		if i.(*int) == nil {
			return "-"
		}
	case "*int64":
		if i.(*int64) == nil {
			return "-"
		}
	}

	return i
}
