package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type logger interface {
	log() *zerolog.Event
	err() *zerolog.Event
	OuteputLog(payload LogPayload)
}

type Logger struct {
	internalLogger zerolog.Logger
}

type LogPayload struct {
	Info  string
	Error error
}

func NewLogger() *Logger {
	lgr := log.Logger.With().Logger()
	logg := &Logger{
		internalLogger: lgr,
	}
	return logg
}

func (l *Logger) log() *zerolog.Event {
	return l.internalLogger.Info()
}

func (l *Logger) err() *zerolog.Event {
	return l.internalLogger.Error()
}

func (l *Logger) OuteputLog(payload LogPayload) {
	if payload.Error != nil {
		l.err().Msg(payload.Error.Error())
	}
	if payload.Info != "" {
		l.log().Msg(payload.Info)
	}
}
