package rtccamlog

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"runtime"
)

func getCallerName() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()
}

func Info() *zerolog.Event {
	return log.Info().Caller(1).Str("Caller", getCallerName())
}

func Error() *zerolog.Event {
	return log.Error().Caller(1).Str("Caller", getCallerName())
}
