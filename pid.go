package main

import (
	"io/ioutil"
	"log"
	"os"
)

var _ = log.Print // for debugging, remove

type PidLib struct {
	path string
}

func NewPidLib(path string) PidLib {
	os.MkdirAll(path, 0700)

	return PidLib{path: path}
}

func (l *PidLib) hasPid(svcName string) bool {
	if _, err := os.Stat(l.path + "/" + svcName); err == nil {
		return true
	}
	return false
}

func (l *PidLib) getPid(svcName string) (string, error) {
	b, err := ioutil.ReadFile(l.path + "/" + svcName)
	if err != nil {
		return "", err
	}

	return string(b[:]), nil
}

func (l *PidLib) setPid(svcName string, id string) error {
	b := []byte(id)

	err := ioutil.WriteFile(l.path+"/"+svcName, b, 0644)

	return err
}
func (l *PidLib) removePid(svcName string) error {
	return os.Remove(l.path + "/" + svcName)
}
