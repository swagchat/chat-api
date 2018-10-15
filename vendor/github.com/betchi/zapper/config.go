package zapper

// Config is settings of logger
type Config struct {
	// EnableConsole is a flag for enable console log.
	EnableConsole bool `yaml:"enableConsole"`
	// ConsoleFormat is a format for console log.
	ConsoleFormat string `yaml:"consoleFormat"`
	// ConsoleLevel is a level for console log.
	ConsoleLevel string `yaml:"consoleLevel"`
	// EnableFile is a flag for enable file log.
	EnableFile bool `yaml:"enableFile"`
	// FileFormat is a format for file log.
	FileFormat string `yaml:"fileFormat"`
	// FileLevel is a log level for file log.
	FileLevel string `yaml:"fileLevel"`
	// FilePath is a file path for file log.
	FilePath string `yaml:"filePath"`
}
