package zapper

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogging(t *testing.T) {
	testCases := []struct {
		Description         string
		LoggerConfiguration *Config
		ExpectedLogs        []string
	}{
		{
			"file logging, json, debug",
			&Config{
				EnableConsole: false,
				EnableFile:    true,
				FileFormat:    "json",
				FileLevel:     levelDebug,
			},
			[]string{
				`{"level":"debug","timestamp":0,"caller":"testing/testing.go:0","msg":"real debug log"}`,
				`{"level":"info","timestamp":0,"caller":"testing/testing.go:0","msg":"real info log"}`,
				`{"level":"warn","timestamp":0,"caller":"testing/testing.go:0","msg":"real warning log"}`,
				`{"level":"error","timestamp":0,"caller":"testing/testing.go:0","msg":"real error log","stacktrace":"dummystacktrace"}`,
			},
		},
		{
			"file logging, json, error",
			&Config{
				EnableConsole: false,
				EnableFile:    true,
				FileFormat:    "json",
				FileLevel:     levelError,
			},
			[]string{
				`{"level":"error","timestamp":0,"caller":"testing/testing.go:0","msg":"real error log","stacktrace":"dummystacktrace"}`,
			},
		},
		{
			"file logging, non-json, debug",
			&Config{
				EnableConsole: false,
				EnableFile:    true,
				FileFormat:    "text",
				FileLevel:     levelDebug,
			},
			[]string{
				`TIME	debug	testing/testing.go:0	real debug log`,
				`TIME	info	testing/testing.go:0	real info log`,
				`TIME	warn	testing/testing.go:0	real warning log`,
				`TIME	error	testing/testing.go:0	real error log`,
			},
		},
		{
			"file logging, non-json, error",
			&Config{
				EnableConsole: false,
				EnableFile:    true,
				FileFormat:    "",
				FileLevel:     levelError,
			},
			[]string{
				`TIME	error	testing/testing.go:0	real error log`,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Description, func(t *testing.T) {
			var filePath string
			if testCase.LoggerConfiguration.EnableFile {
				tempDir, err := ioutil.TempDir(os.TempDir(), "TestLoggingAfterInitialized")
				require.NoError(t, err)
				defer os.Remove(tempDir)

				filePath = filepath.Join(tempDir, "file.log")
				testCase.LoggerConfiguration.FilePath = filePath
			}

			logger := NewLogger(testCase.LoggerConfiguration)

			logger.Debug("real debug log")
			logger.Info("real info log")
			logger.Warn("real warning log")
			logger.Error("real error log")

			if testCase.LoggerConfiguration.EnableFile {
				logs, err := ioutil.ReadFile(filePath)
				require.NoError(t, err)

				actual := strings.TrimSpace(string(logs))

				if testCase.LoggerConfiguration.FileFormat == "json" {
					reTs := regexp.MustCompile(`"timestamp":"([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2}):([0-9]{2})\.([0-9]{3})\+([0-9]{4})"`)
					reCaller := regexp.MustCompile(`"caller":"([^"]+):[0-9\.]+"`)
					reStacktrace := regexp.MustCompile(`"stacktrace":".*"`)
					actual = reTs.ReplaceAllString(actual, `"timestamp":0`)
					actual = reCaller.ReplaceAllString(actual, `"caller":"$1:0"`)
					actual = reStacktrace.ReplaceAllString(actual, `"stacktrace":"dummystacktrace"`)
				} else {
					actualRows := strings.Split(actual, "\n")
					for i, actualRow := range actualRows {
						actualFields := strings.Split(actualRow, "\t")
						if len(actualFields) > 3 {
							actualFields[0] = "TIME"
							reCaller := regexp.MustCompile(`([^"]+):[0-9\.]+`)
							actualFields[2] = reCaller.ReplaceAllString(actualFields[2], "$1:0")
							actualRows[i] = strings.Join(actualFields, "\t")
						}
					}

					actual = strings.Join(actualRows, "\n")
				}
				require.Equal(t, testCase.ExpectedLogs, strings.Split(actual, "\n"))
			}
		})
	}
}
