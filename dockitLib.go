package main

import (
	"errors"
	"fmt"
	"github.com/dotcloud/docker"
	dockerc "github.com/fsouza/go-dockerclient"
	"log"
)

var _ = log.Print // for debugging, remove

type Lib struct {
	cfg     ServiceMap
	address string
	client  *dockerc.Client
	pids    PidLib
}

func NewLib(cfg map[string]Service, address string, pidPath string) Lib {
	c, err := dockerc.NewClient(address)
	if err != nil {
		panic(err)
	}
	pids := NewPidLib(pidPath)
	return Lib{cfg: cfg, address: address, client: c, pids: pids}
}

func (l *Lib) Start(svcName string) error {
	image := l.cfg[svcName].Image
	ports := l.cfg[svcName].Ports
	env := l.cfg[svcName].Env

	if l.pids.hasPid(svcName) {
		return errors.New("Service " + svcName + " already running")
	}

	// Start Dependency Containers
	if err := l.startDeps(svcName); err != nil {
		return err
	}

	// Create Container
	opts := dockerc.CreateContainerOptions{}
	config := docker.Config{
		Image: image,
		Env:   l.getEnv(env),
	}
	container, err := l.client.CreateContainer(opts, &config)
	if err != nil {
		return err
	}

	// Start Container
	links, err := l.getLinks(svcName)
	if err != nil {
		return err
	}
	hostConfig := docker.HostConfig{
		PortBindings: l.getPortBindings(ports),
		Links:        links,
	}
	err = l.client.StartContainer(container.ID, &hostConfig)
	if err != nil {
		return err
	}
	if err = l.pids.setPid(svcName, container.ID); err != nil {
		return err
	}
	return nil
}

func (l *Lib) Stop(svcName string) error {
	if !l.pids.hasPid(svcName) {
		return errors.New("Service not running")
	}

	id, err := l.pids.getPid(svcName)
	if err != nil {
		return err
	}

	if err = l.client.StopContainer(id, 5); err != nil {
		return err
	}

	if err = l.pids.removePid(svcName); err != nil {
		return err
	}
	return nil
}
func (l *Lib) startDeps(svcName string) error {
	deps := l.cfg[svcName].Deps

	for _, svcName := range deps {
		if !l.pids.hasPid(svcName) {

			if err := l.Start(svcName); err != nil {
				return err
			}
			fmt.Println("Dep " + svcName + " started")
		}
	}
	return nil
}
func (l *Lib) getContainerName(svcName string) (string, error) {
	c, err := dockerc.NewClient(l.address)
	if err != nil {
		return "", err
	}
	id, err := l.pids.getPid(svcName)
	if err != nil {
		return "", err
	}
	container, err := c.InspectContainer(id)
	if err != nil {
		return "", err
	}
	return container.Name, nil
}

func (l *Lib) getPortBindings(ports map[string]string) map[docker.Port][]docker.PortBinding {
	portBindings := make(map[docker.Port][]docker.PortBinding)

	for internal, external := range ports {
		portBinding := []docker.PortBinding{}
		if external != "" {
			portBinding = []docker.PortBinding{docker.PortBinding{HostIp: "0.0.0.0", HostPort: external}}
		}
		port := docker.NewPort("tcp", internal)

		portBindings[port] = portBinding
	}

	return portBindings
}
func (l *Lib) getEnv(env map[string]string) []string {
	envFlat := make([]string, 0, 10)
	for k, v := range env {
		envFlat = append(envFlat, k+"="+v)
	}
	return envFlat
}
func (l *Lib) getLinks(svcName string) ([]string, error) {
	deps := l.cfg[svcName].Deps
	links := make([]string, 0, 10)
	for _, svcName := range deps {
		if !l.pids.hasPid(svcName) {
			return links, errors.New("Dep not running: " + svcName)
		}
		name, err := l.getContainerName(svcName)
		if err != nil {
			return links, err
		}
		links = append(links, name+":"+svcName)
	}
	return links, nil
}
