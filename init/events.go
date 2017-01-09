// Copyright 2017 Saad Nasser (SdNssr). All rights reserved.
//
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package main

import (
	"os"
	"os/signal"
	"syscall"
)

var SIGMIN = 1
var SIGMAX = 31
var SIGRTMIN = 34
var SIGRTMAX = 64

var PEACEFUL_HALT = syscall.Signal(SIGRTMIN + 3)
var PEACEFUL_SHUTDOWN = syscall.Signal(SIGRTMIN + 4)
var PEACEFUL_REBOOT = syscall.Signal(SIGRTMIN + 5)

var FORCEFUL_HALT = syscall.Signal(SIGRTMIN + 13)
var FORCEFUL_SHUTDOWN = syscall.Signal(SIGRTMIN + 14)
var FORCEFUL_REBOOT = syscall.Signal(SIGRTMIN + 15)

var RECOVERY = syscall.Signal(SIGRTMIN + 2)
var RESCUE = syscall.Signal(SIGRTMIN + 1)
var NORMAL = syscall.Signal(SIGRTMIN + 0)

var POWERFAIL = syscall.SIGPWR
var KBREQ = syscall.SIGWINCH
var SAK = syscall.SIGINT

var REAP = syscall.SIGCHLD

func NotifyAllSignals(c chan<- os.Signal) {
	for i := SIGMIN; i <= SIGMAX; i++ {
		signal.Notify(c, syscall.Signal(i))
	}

	for i := SIGRTMIN; i <= SIGRTMAX; i++ {
		signal.Notify(c, syscall.Signal(i))
	}
}
