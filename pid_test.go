package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var _ = fmt.Print // For debugging; delete when done.
var _ = log.Print // For debugging; delete when done.
var _ = os.Stdout // For debugging; delete when done.

var tmpPath = "/tmp/pid-tests"
var testSvc = "ServiceName"
var testId = "jhgsd765asd"

func Test_PidLib_hasPid(t *testing.T) {
	// given
	pids := NewPidLib(tmpPath + "a")
	defer os.RemoveAll(tmpPath + "a")
	pids.setPid(testSvc, testId)

	// when
	// then
	if !pids.hasPid(testSvc) {
		t.Errorf("pid should exist")
	}
}

func Test_PidLib_getPid(t *testing.T) {
	// given
	pids := NewPidLib(tmpPath + "b")
	defer os.RemoveAll(tmpPath + "b")
	pids.setPid(testSvc, testId)

	// when
	// then
	if id, _ := pids.getPid(testSvc); id != testId {
		t.Errorf("pid should equal " + testId)
	}
}

func Test_PidLib_removePid(t *testing.T) {
	// given
	pids := NewPidLib(tmpPath + "c")
	defer os.RemoveAll(tmpPath + "c")
	pids.setPid(testSvc, testId)

	// when
	pids.removePid(testSvc)

	// then
	if pids.hasPid(testSvc) {
		t.Errorf("pid should NOT exist")
	}
}
