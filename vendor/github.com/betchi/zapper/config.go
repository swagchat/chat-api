package zapper

// Config is settings of logger
type Config struct {
	// EnableConsole is a flag for enable console log.
	EnableConsole bool
	// ConsoleFormat is a format for console log.
	ConsoleFormat string
	// ConsoleLevel is a level for console log.
	ConsoleLevel string

	// EnableFile is a flag for enable file log.
	EnableFile bool
	// FileFormat is a format for file log.
	FileFormat string
	// FileLevel is a log level for file log.
	FileLevel string
	// FilePath is a file path for file log.
	FilePath string
	// FileMaxSize is a max size(megabytes) for file log. The Default is 100 megabytes.
	FileMaxSize int
	// FileMaxAge is a max age(days) for file log. The default is not to remove old log files.
	FileMaxAge int
	// FileMaxBackups is a max backup file number for file log. The default is to retain all old log files
	FileMaxBackups int
	// FileLocalTime is a formatting the timestamps for backup files. The default is to use UTC time.
	FileLocalTime bool
	// FileCompress is a compressed flag for file log. The default is false.
	FileCompress bool
}
