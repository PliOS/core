// Copyright 2017 Saad Nasser (SdNssr). All rights reserved.
//
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package main

import (
	log "github.com/Sirupsen/logrus"
	"os/exec"
	"os"
)

type ServiceAction struct {
	Service string
	Action  string
}

type ExecCommand struct {
	Program   string
	Arguments []string
}

func RunServices(pidDied chan int, serviceActions chan ServiceAction, execCommands chan ExecCommand, execFinished chan int) {
	execPid := 0

	for {
		select {
		case processDied := <- pidDied:
			log.WithFields(log.Fields{
				"pid": processDied,
			}).Debugf("Process died")

			if processDied == execPid {
				execFinished <- execPid
			}
		case serviceAction := <- serviceActions:
			log.WithFields(log.Fields{
				"service": serviceAction.Service,
				"action": serviceAction.Action,
			}).Debugf("Service action")
		case execCommand := <- execCommands:
			log.WithFields(log.Fields{
				"program": execCommand.Program,
				"arguments": execCommand.Arguments,
			}).Debugf("Executing command")


			cmd := exec.Command(execCommand.Program, execCommand.Arguments...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Start(); err != nil {
				log.Fatalf("Fatal error - exec(%s, %v): %s", execCommand.Program, execCommand.Arguments, err)
			}

			execPid = cmd.Process.Pid
		}
	}
}
