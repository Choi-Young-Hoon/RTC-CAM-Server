package rtccamlog

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Info() *zerolog.Event {
	return log.Info().Caller(1)
}

func Error() *zerolog.Event {
	return log.Error().Caller()
}
