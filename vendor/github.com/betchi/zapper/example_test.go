package zapper_test

import logger "github.com/betchi/zapper"

// ExampleConsoleTextTest logs text format to console
func ExampleConsoleText() {
	logger.InitGlobalLogger(&logger.Config{
		EnableConsole: true,
		ConsoleFormat: "text",
		ConsoleLevel:  "debug",
	})

	logger.Debug("debug level")
	logger.Info("info level")
	logger.Warn("warn level")
	logger.Error("error level")
	logger.Fatal("fatal level")
}

// ExampleConsoleJson logs json format to console
func ExampleConsoleJson() {
	logger.InitGlobalLogger(&logger.Config{
		EnableConsole: true,
		ConsoleFormat: "json",
		ConsoleLevel:  "debug",
	})

	logger.Debug("debug level")
	logger.Info("info level")
	logger.Warn("warn level")
	logger.Error("error level")
	logger.Fatal("fatal level")
}

// ExampleFileText logs text format to file
func ExampleFileText() {
	logger.InitGlobalLogger(&logger.Config{
		EnableFile: true,
		FileFormat: "text",
		FileLevel:  "debug",
		FilePath:   "app.log",
	})

	logger.Debug("debug level")
	logger.Info("info level")
	logger.Warn("warn level")
	logger.Error("error level")
	logger.Fatal("fatal level")
}

// ExampleConsoleAndFileJson logs json format to both console and file
func ExampleConsoleAndFileJson() {
	logger.InitGlobalLogger(&logger.Config{
		EnableConsole: true,
		ConsoleFormat: "json",
		ConsoleLevel:  "debug",
		EnableFile:    true,
		FileFormat:    "json",
		FileLevel:     "debug",
		FilePath:      "app.log",
	})

	logger.Debug("debug level")
	logger.Info("info level")
	logger.Warn("warn level")
	logger.Error("error level")
	logger.Fatal("fatal level")
}
