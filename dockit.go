package main

import (
	"flag"
	"fmt"
	"log"
)

var _ = log.Print

func main() {

	config := flag.String("config", "config.json", "json config")
	address := flag.String("address", "unix:///var/run/docker.sock", "docker address")
	pidPath := flag.String("pidPath", "/var/run/dockit-containers", "path to store pids in")
	service := flag.String("service", "", "service name")
	start := flag.Bool("start", false, "start service")
	stop := flag.Bool("stop", false, "stop service")
	flag.Parse()

	cfg, err := parseConfigFile(*config)
	if err != nil {
		panic(err)
	}
	if *service == "" {
		flag.Usage()
		return
	}

	lib := NewLib(cfg, *address, *pidPath)

	svcName := *service

	switch {
	case *start:
		if err = lib.Start(svcName); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(svcName + " Started")
	case *stop:
		if err = lib.Stop(svcName); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(svcName + " Stopped")
	default:
		flag.Usage()
	}
}
