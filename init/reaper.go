// Copyright 2017 Saad Nasser (SdNssr). All rights reserved.
//
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

package main

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
)

func ReapChildren(pids chan<- int) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, unix.SIGCHLD)

	for range signals {
		for {
			var status unix.WaitStatus
			pid, err := unix.Wait4(-1, &status, unix.WNOHANG, nil)

			switch err {
			case nil:
				if pid <= 0 {
					break
				} else {
					log.WithFields(log.Fields{
						"pid": pid,
					}).Debugf("Reaped process")

					pids <- pid
				}
			default:
				break
			}
		}
	}
}
