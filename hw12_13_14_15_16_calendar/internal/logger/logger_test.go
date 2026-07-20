package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger_DebugFilteredAtInfoLevel(t *testing.T) {
	var buf bytes.Buffer
	log := New(Conf{Level: InfoLevel}, &buf)

	log.Info("info msg")
	log.Debug("debug msg")

	assert.Contains(t, buf.String(), "info msg")
	assert.NotContains(t, buf.String(), "debug msg")
}

func TestLogger_ErrorAlwaysWrittenAtErrorLevel(t *testing.T) {
	var buf bytes.Buffer
	log := New(Conf{Level: ErrorLevel}, &buf)

	log.Error("error msg")
	log.Info("info msg")

	assert.Contains(t, buf.String(), "error msg")
	assert.NotContains(t, buf.String(), "info msg")
}
