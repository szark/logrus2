package logrus2

// Portions of this files were taken from https://github.com/szark/logrus2
// Copyrighted by Simon Eskildsen under the MIT license

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

// Defines a log format type that wil output line separated JSON objects
// in the GELF format.
type GelfFormatter struct {
}

type fields map[string]interface{}

var (
	DefaultLevel = LOG_INFO
	levelMap     = map[Level]Priority{
		PanicLevel: LOG_EMERG,
		FatalLevel: LOG_CRIT,
		ErrorLevel: LOG_ERR,
		WarnLevel:  LOG_WARNING,
		InfoLevel:  LOG_INFO,
		DebugLevel: LOG_DEBUG,
		TraceLevel: LOG_TRACE,
	}
)

// Format formats the log entry to GELF JSON
func (f *GelfFormatter) Format(entry *Entry) ([]byte, error) {
	data := make(fields, len(entry.Data)+6)
	blacklist := []string{"_id", "id", "timestamp", "version", "level"}

	for k, v := range entry.Data {

		if contains(k, blacklist) {
			continue
		}

		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			data["_"+k] = v.Error()
		default:
			data["_"+k] = v
		}
	}

	data["version"] = "1.1"
	data["short_message"] = entry.Message
	data["timestamp"] = round((float64(entry.Time.UnixNano())/float64(1000000))/float64(1000), 4)
	data["level"] = toSyslogLevel(entry.Level)
	data["level_name"] = entry.Level.String()
	data["_pid"] = os.Getpid()

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	fmt.Println("Gelf formatter format func?")
	return append(serialized, '\n'), nil
}

func contains(needle string, haystack []string) bool {
	for _, a := range haystack {
		if needle == a {
			return true
		}
	}
	return false
}

func round(val float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Floor((val*shift)+.5) / shift
}

func toSyslogLevel(level Level) Priority {
	syslog, ok := levelMap[level]
	if ok {
		return syslog
	}
	return DefaultLevel
}
