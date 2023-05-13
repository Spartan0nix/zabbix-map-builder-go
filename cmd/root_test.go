package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var mappingFilePath string

func init() {
	pwd, _ := os.Getwd()
	mappingFilePath = filepath.Join(pwd, "..", "examples", "mapping.json")
}

func TestCheckRequiredFlag(t *testing.T) {
	checkRequiredFlag("test", mappingFilePath)
}

func TestCheckRequiredFlagMissingName(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		checkRequiredFlag("", mappingFilePath)

		return
	}

	// Execute test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestCheckRequiredFlagMissingName")
	// Run the desired command when running test in suprocess
	cmd.Env = append(cmd.Env, "BE_CRASHER=1")
	err := cmd.Run()

	exit, _ := err.(*exec.ExitError)
	if exit.ExitCode() != 1 {
		t.Fatalf("expected exit code 1.\nCode returned : %d\nError returned : %s", exit.ExitCode(), err)
	}
}

func TestCheckRequiredFlagMissingFile(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		checkRequiredFlag("test", "")

		return
	}

	// Execute test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestCheckRequiredFlagMissingName")
	// Run the desired command when running test in suprocess
	cmd.Env = append(cmd.Env, "BE_CRASHER=1")
	err := cmd.Run()

	exit, _ := err.(*exec.ExitError)
	if exit.ExitCode() != 1 {
		t.Fatalf("expected exit code 1.\nCode returned : %d\nError returned : %s", exit.ExitCode(), err)
	}
}
