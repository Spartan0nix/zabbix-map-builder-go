package logging

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

var testLogger *Logger

func init() {
	testLogger = NewLogger(Debug)
}

func TestNewLogger(t *testing.T) {
	l := NewLogger(Debug)
	if l == nil {
		t.Fatal("A nil pointer was returned instead of *Logger")
	}
}

func BenchmarkNewLogger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewLogger(Debug)
	}
}

func TestGetLevel(t *testing.T) {
	if level := getLevel(Critical); level != "CRITICAL" {
		t.Fatalf("Wrong string representation returned\nExpected '%s'\nReturned : %s", "CRITICAL", level)
	}

	if level := getLevel(Error); level != "ERROR" {
		t.Fatalf("Wrong string representation returned\nExpected '%s'\nReturned : %s", "ERROR", level)
	}

	if level := getLevel(Warning); level != "WARNING" {
		t.Fatalf("Wrong string representation returned\nExpected '%s'\nReturned : %s", "WARNING", level)
	}

	if level := getLevel(Info); level != "INFO" {
		t.Fatalf("Wrong string representation returned\nExpected '%s'\nReturned : %s", "INFO", level)
	}

	if level := getLevel(Debug); level != "DEBUG" {
		t.Fatalf("Wrong string representation returned\nExpected '%s'\nReturned : %s", "DEBUG", level)
	}

}

func BenchmarkGetLevel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getLevel(Debug)
	}
}

func TestSetFlags(t *testing.T) {
	testLogger.setFlags()
	if testLogger.logger.Flags() != log.Lmsgprefix {
		t.Fatalf("Wrong logger flags set\nExpected : %d\nReturned : %d", log.Lmsgprefix, testLogger.logger.Flags())
	}
}

func BenchmarkSetFlags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testLogger.setFlags()
	}
}

func TestWriteLog(t *testing.T) {
	var buf bytes.Buffer
	expectedOutput := "[map-builder][INFO] test-value\n"

	// Set the logger output to a buffer instead of os.Stderr file
	testLogger.logger.SetOutput(&buf)

	testLogger.writeLog(Info, "test-value")

	bufOut := buf.String()
	if bufOut != expectedOutput {
		t.Fatalf("Wrong log format returned\nExpected : %s\nReturned : %s", expectedOutput, bufOut)
	}
}

func TestWriteLogPanic(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("panic function was not trigger")
		}
	}()

	var buf bytes.Buffer
	expectedOutput := "[map-builder][CRITICAL] test-value\n"

	// Set the logger output to a buffer instead of os.Stderr file
	testLogger.logger.SetOutput(&buf)

	testLogger.writeLog(Critical, "test-value")

	if !strings.Contains(buf.String(), expectedOutput) {
		t.Fatalf("Wrong log format returned\nExpected '%s' to be present in the output string\nReturned : %s", expectedOutput, buf.String())
	}
}

func BenchmarkWriteLog(b *testing.B) {
	var buf bytes.Buffer
	testLogger.logger.SetOutput(&buf)

	for i := 0; i < b.N; i++ {
		testLogger.writeLog(Info, "bench-value")
	}
}

func TestCritical(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("panic function was not trigger")
		}
	}()

	var buf bytes.Buffer
	expectedOutput := "[map-builder][CRITICAL] test-value\n"

	// Set the logger output to a buffer instead of os.Stderr file
	testLogger.logger.SetOutput(&buf)

	testLogger.Critical("test-value")

	if !strings.Contains(buf.String(), expectedOutput) {
		t.Fatalf("Wrong log format returned\nExpected '%s' to be present in the output string\nReturned : %s", expectedOutput, buf.String())
	}
}

func TestError(t *testing.T) {
	var buf bytes.Buffer
	expectedOutput := "[map-builder][ERROR] test-value\n"

	// Set the logger output to a buffer instead of os.Stderr file
	testLogger.logger.SetOutput(&buf)

	testLogger.Error("test-value")

	if !strings.Contains(buf.String(), expectedOutput) {
		t.Fatalf("Wrong log format returned\nExpected '%s' to be present in the output string\nReturned : %s", expectedOutput, buf.String())
	}
}

func TestWarning(t *testing.T) {
	var buf bytes.Buffer
	expectedOutput := "[map-builder][WARNING] test-value\n"

	// Set the logger output to a buffer instead of os.Stderr file
	testLogger.logger.SetOutput(&buf)

	testLogger.Warning("test-value")

	if buf.String() != expectedOutput {
		t.Fatalf("Wrong log format returned\nExpected : %s\nReturned : %s", expectedOutput, buf.String())
	}
}

func TestInfo(t *testing.T) {
	var buf bytes.Buffer
	expectedOutput := "[map-builder][INFO] test-value\n"

	// Set the logger output to a buffer instead of os.Stderr file
	testLogger.logger.SetOutput(&buf)

	testLogger.Info("test-value")

	if buf.String() != expectedOutput {
		t.Fatalf("Wrong log format returned\nExpected : %s\nReturned : %s", expectedOutput, buf.String())
	}
}

func TestDebug(t *testing.T) {
	var buf bytes.Buffer
	expectedOutput := "[map-builder][DEBUG] test-value\n"

	// Set the logger output to a buffer instead of os.Stderr file
	testLogger.logger.SetOutput(&buf)

	testLogger.Debug("test-value")

	if buf.String() != expectedOutput {
		t.Fatalf("Wrong log format returned\nExpected : %s\nReturned : %s", expectedOutput, buf.String())
	}
}
