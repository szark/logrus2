package logrus2

// Copied from log/syslog as that package is not available on windows.

type Priority int

const (
	// Severity.

	// From /usr/include/sys/syslog.h.
	// These are the same on Linux, BSD, and OS X.
	LOG_EMERG Priority = iota
	LOG_CRIT
	LOG_ERR
	LOG_WARNING
	LOG_INFO
	LOG_DEBUG
	LOG_TRACE
)
