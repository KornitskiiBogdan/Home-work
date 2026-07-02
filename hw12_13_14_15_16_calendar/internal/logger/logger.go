package logger

import (
	"io"

	"github.com/rs/zerolog"
)

type Level string

const (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
)

type Conf struct {
	Level Level `yaml:"level"`
}

type logger struct {
	log *zerolog.Logger
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg string)
}

func New(conf Conf, writer io.Writer) Logger {
	var levelZerolog zerolog.Level
	switch conf.Level {
	case DebugLevel:
		levelZerolog = zerolog.DebugLevel
	case InfoLevel:
		levelZerolog = zerolog.InfoLevel
	case WarnLevel:
		levelZerolog = zerolog.WarnLevel
	case ErrorLevel:
		levelZerolog = zerolog.ErrorLevel
	default:
		levelZerolog = zerolog.DebugLevel
	}

	zeroLog := zerolog.New(writer).
		With().
		Timestamp().
		Logger().
		Level(levelZerolog)

	return &logger{
		log: &zeroLog,
	}
}

func (l *logger) Info(msg string) {
	l.log.Info().Msg(msg)
}

func (l *logger) Error(msg string) {
	l.log.Error().Msg(msg)
}

func (l *logger) Warn(msg string) {
	l.log.Warn().Msg(msg)
}

func (l *logger) Debug(msg string) {
	l.log.Debug().Msg(msg)
}
