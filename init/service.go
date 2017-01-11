// Copyright 2017 Saad Nasser (SdNssr). All rights reserved.
//
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package main

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
)

type ServiceAction struct {
	Service string
	Action  string
}

type ExecCommand struct {
	Program   string
	Arguments []string
}

func RunCommand(program string, arguments []string) int {
	cmd := exec.Command(program, arguments...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("Fatal error - exec(%s, %v): %s", program, arguments, err)
	}

	return cmd.Process.Pid
}

func RunServices(config *Config, pidDied chan int, serviceActions chan ServiceAction, execCommands chan ExecCommand, execFinished chan int) {
	execPid := 0

	runningServices := map[string]int{}
	servicePids := map[int]string{}

	for {
		select {
		case processDied := <-pidDied:
			log.WithFields(log.Fields{
				"pid": processDied,
			}).Debugf("Process died")

			if processDied == execPid {
				execFinished <- execPid
			}

			if service, running := servicePids[processDied]; running {
				if _, active := runningServices[service]; active {
					runningServices[service] = RunCommand(
						config.Services[service].Program,
						config.Services[service].Arguments,
					)

					servicePids[runningServices[service]] = service

					log.WithFields(log.Fields{
						"name":   service,
						"oldPid": processDied,
						"newPid": runningServices[service],
					}).Debugf("Respawned service")
				} else {
					delete(servicePids, processDied)

					log.WithFields(log.Fields{
						"name": service,
						"pid":  processDied,
					}).Debugf("Service died")
				}
			}
		case serviceAction := <-serviceActions:
			log.WithFields(log.Fields{
				"service": serviceAction.Service,
				"action":  serviceAction.Action,
			}).Debugf("Service action")

			switch serviceAction.Action {
			case "start":
				if _, running := runningServices[serviceAction.Service]; !running {
					runningServices[serviceAction.Service] = RunCommand(
						config.Services[serviceAction.Service].Program,
						config.Services[serviceAction.Service].Arguments,
					)

					servicePids[runningServices[serviceAction.Service]] = serviceAction.Service

					log.WithFields(log.Fields{
						"name": serviceAction.Service,
						"pid":  runningServices[serviceAction.Service],
					}).Debugf("Started service")
				}
			case "stop":
				if pid, running := runningServices[serviceAction.Service]; running {
					unix.Kill(pid, unix.SIGTERM)
					delete(runningServices, serviceAction.Service)

					log.WithFields(log.Fields{
						"name": serviceAction.Service,
						"pid":  pid,
					}).Debugf("Stopped service")
				}
			case "restart":
				if pid, running := runningServices[serviceAction.Service]; running {
					unix.Kill(pid, unix.SIGTERM)

					log.WithFields(log.Fields{
						"name": serviceAction.Service,
						"pid":  pid,
					}).Debugf("Restarted service")
				}
			}
		case execCommand := <-execCommands:
			log.WithFields(log.Fields{
				"program":   execCommand.Program,
				"arguments": execCommand.Arguments,
			}).Debugf("Executing command")

			execPid = RunCommand(execCommand.Program, execCommand.Arguments)
		}
	}
}
