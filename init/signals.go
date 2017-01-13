// Copyright (c) 2017 The pliOS Authors. All rights reserved.
//
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var SIGMIN = 1
var SIGMAX = 31
var SIGRTMIN = 34
var SIGRTMAX = 64

var SIGNAL_RECOVERY = syscall.Signal(SIGRTMIN + 2)
var SIGNAL_HALT = syscall.Signal(SIGRTMIN + 3)
var SIGNAL_SHUTDOWN = syscall.Signal(SIGRTMIN + 4)
var SIGNAL_REBOOT = syscall.Signal(SIGRTMIN + 5)

func NotifyAllSignals(c chan<- os.Signal) {
	for i := SIGMIN; i <= SIGMAX; i++ {
		signal.Notify(c, syscall.Signal(i))
	}

	for i := SIGRTMIN; i <= SIGRTMAX; i++ {
		signal.Notify(c, syscall.Signal(i))
	}
}

func ProcessSignals(triggerRunner *TriggerRunner) {
	signals := make(chan os.Signal, 1)
	NotifyAllSignals(signals)

	for signal := range signals {
		log.WithFields(log.Fields{
			"signal": signal,
		}).Debugf("Received signal")

		switch signal {
		case SIGNAL_HALT:
			triggerRunner.RunTrigger("halt")
		case SIGNAL_SHUTDOWN:
			triggerRunner.RunTrigger("shutdown")
		case SIGNAL_REBOOT:
			triggerRunner.RunTrigger("reboot")
		case SIGNAL_RECOVERY:
			triggerRunner.RunTrigger("recovery")
		}
	}
}
