// Copyright 2017 Saad Nasser (SdNssr). All rights reserved.
//
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stderr)

	SetupProcessEnvironment()
	MountApiFilesystems()
	CreateApiSymlinks()

	log.Infof("Initalized system")

	config := ReadConfig()

	log.Infof("Read config file")

	serviceActions := make(chan ServiceAction)
	execCommands := make(chan ExecCommand)
	execFinished := make(chan int)
	triggers := make(chan string)
	pids := make(chan int)

	go func() { triggers <- "init" }()

	go ReapChildren(pids)
	go ProcessSignals(triggers)
	go RunServices(config, pids, serviceActions, execCommands, execFinished)
	ProcessTriggers(config, triggers, execCommands, execFinished, serviceActions)
}
