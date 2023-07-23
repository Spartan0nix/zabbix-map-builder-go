package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

const (
	testFile = "test.json"
)

var mappingFilePath string

func init() {
	pwd, _ := os.Getwd()
	mappingFilePath = filepath.Join(pwd, "..", "..", "examples", "mapping.json")
}

func TestGetEnvironmentVariables(t *testing.T) {
	if err := os.Setenv("ZABBIX_URL", ZABBIX_URL); err != nil {
		t.Fatalf("error while setting environment variable 'ZABBIX_URL'.\nReason : %v", err)
	}

	if err := os.Setenv("ZABBIX_USER", ZABBIX_USER); err != nil {
		t.Fatalf("error while setting environment variable 'ZABBIX_USER'.\nReason : %v", err)
	}

	if err := os.Setenv("ZABBIX_PWD", ZABBIX_PWD); err != nil {
		t.Fatalf("error while setting environment variable 'ZABBIX_PWD' .\nReason : %v", err)
	}

	opts, err := GetEnvironmentVariables()
	if err != nil {
		t.Fatalf("error while executing GetEnvironmentVariables function.\nReason : %v", err)
	}

	if opts.ZabbixUrl != ZABBIX_URL {
		t.Fatalf("wrong value assigned to 'ZabbixUrl'.\nExpected : %v\nReturned : %v", ZABBIX_URL, opts.ZabbixUrl)
	}
	if opts.ZabbixUser != ZABBIX_USER {
		t.Fatalf("wrong value assigned to 'ZabbixUser'.\nExpected : %v\nReturned : %v", ZABBIX_USER, opts.ZabbixUser)
	}
	if opts.ZabbixPwd != ZABBIX_PWD {
		t.Fatalf("wrong value assigned to 'ZabbixPwd'.\nExpected : %v\nReturned : %v", ZABBIX_PWD, opts.ZabbixPwd)
	}
}

func TestGetEnvironmentVariablesMissingUrl(t *testing.T) {
	if err := os.Unsetenv("ZABBIX_URL"); err != nil {
		t.Fatalf("error while unsetting environment variable 'ZABBIX_URL'.\nReason : %v", err)
	}

	if err := os.Setenv("ZABBIX_USER", ZABBIX_USER); err != nil {
		t.Fatalf("error while setting environment variable 'ZABBIX_USER'.\nReason : %v", err)
	}

	if err := os.Setenv("ZABBIX_PWD", ZABBIX_PWD); err != nil {
		t.Fatalf("error while setting environment variable 'ZABBIX_PWD' .\nReason : %v", err)
	}

	opts, err := GetEnvironmentVariables()
	if err == nil {
		t.Fatalf("an error should be returned when an environment variable is missing")
	}

	if opts != nil {
		t.Fatalf("a nil pointer should be returned instead of *Options when an environment variable is missing")
	}
}

func TestGetEnvironmentVariablesMissingUser(t *testing.T) {
	if err := os.Setenv("ZABBIX_URL", ZABBIX_URL); err != nil {
		t.Fatalf("error while setting environment variable 'ZABBIX_URL'.\nReason : %v", err)
	}

	if err := os.Unsetenv("ZABBIX_USER"); err != nil {
		t.Fatalf("error while unsetting environment variable 'ZABBIX_USER'.\nReason : %v", err)
	}

	if err := os.Setenv("ZABBIX_PWD", ZABBIX_PWD); err != nil {
		t.Fatalf("error while setting environment variable 'ZABBIX_PWD' .\nReason : %v", err)
	}

	opts, err := GetEnvironmentVariables()
	if err == nil {
		t.Fatalf("an error should be returned when an environment variable is missing")
	}

	if opts != nil {
		t.Fatalf("a nil pointer should be returned instead of *Options when an environment variable is missing")
	}
}

func TestGetEnvironmentVariablesMissingPwd(t *testing.T) {
	if err := os.Setenv("ZABBIX_URL", ZABBIX_URL); err != nil {
		t.Fatalf("error while setting environment variable 'ZABBIX_URL'.\nReason : %v", err)
	}

	if err := os.Setenv("ZABBIX_USER", ZABBIX_USER); err != nil {
		t.Fatalf("error while setting environment variable 'ZABBIX_USER'.\nReason : %v", err)
	}

	if err := os.Unsetenv("ZABBIX_PWD"); err != nil {
		t.Fatalf("error while unsetting environment variable 'ZABBIX_PWD'.\nReason : %v", err)
	}

	opts, err := GetEnvironmentVariables()
	if err == nil {
		t.Fatalf("an error should be returned when an environment variable is missing")
	}

	if opts != nil {
		t.Fatalf("a nil pointer should be returned instead of *Options when an environment variable is missing")
	}
}

func TestReadInput(t *testing.T) {
	m, err := ReadInput(mappingFilePath)
	if err != nil {
		t.Fatalf("error while executing ReadInput function.\nReason : %v", err)
	}

	if len(m) == 0 {
		t.Fatal("an empty list of mappings was retrieve")
	}
}

func TestReadInputMissingFile(t *testing.T) {
	m, err := ReadInput(testFile)
	if err == nil {
		t.Fatalf("an error should be returned when the given file does not exist")
	}

	if m != nil {
		t.Fatal("a nil pointer should be returned instead of *[]zbxMap.Mapping when the processing fails")
	}
}

func TestReadInputWrongType(t *testing.T) {
	data := MapOptions{
		ZabbixUrl:  ZABBIX_URL,
		ZabbixUser: ZABBIX_USER,
	}

	b, err := json.Marshal(&data)
	if err != nil {
		t.Fatalf("error while converting data to a slice of byte.\nReason : %v", err)
	}

	err = os.WriteFile(testFile, b, 0644)
	if err != nil {
		t.Fatalf("an error while writing test data to file '%s'.\nReason : %v", testFile, err)
	}

	m, err := ReadInput(testFile)
	if err == nil {
		t.Fatalf("an error should be returned during file processing when json.Unmarshal fail")
	}

	if m != nil {
		t.Fatal("a nil pointer should be returned instead of *[]zbxMap.Mapping when the processing fails")
	}

	err = os.Remove(testFile)
	if err != nil {
		t.Fatalf("error while removing file '%s'.\nReason : %v", testFile, err)
	}
}

func TestReadInputNoData(t *testing.T) {
	err := os.WriteFile(testFile, []byte{}, 0644)
	if err != nil {
		t.Fatalf("an error while writing test data to file '%s'.\nReason : %v", testFile, err)
	}

	m, err := ReadInput(testFile)
	if err == nil {
		t.Fatalf("an error should be returned when no mappings are found in the given file")
	}

	if m != nil {
		t.Fatal("a nil pointer should be returned instead of *[]zbxMap.Mapping when the processing fails")
	}

	err = os.Remove(testFile)
	if err != nil {
		t.Fatalf("error while removing file '%s'.\nReason : %v", testFile, err)
	}
}
