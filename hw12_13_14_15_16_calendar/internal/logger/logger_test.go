package logger

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoggerLevels(t *testing.T) {
	var buf bytes.Buffer

	// подменяем файл и stdout на буфер
	tmpFile, err := os.CreateTemp("", "logger_test")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	multi := io.MultiWriter(&buf, tmpFile)
	l := &Logger{
		level:  INFO,
		logger: log.New(multi, "", 0),
	}

	l.Debug("debug line") // не должен попасть
	l.Info("info line")   // должен попасть
	l.Error("error line") // должен попасть
	output := buf.String()

	require.NotContains(t, output, "debug line")
	require.Contains(t, output, "info line")
	require.Contains(t, output, "error line")
}
