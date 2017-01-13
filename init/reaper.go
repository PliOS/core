// Copyright (c) 2017 The pliOS Authors. All rights reserved.
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

type PidWaitChannel chan unix.WaitStatus
type WildcardWaitChannel chan WaitStatusPid

type WaitStatusPid struct {
	Pid    int
	Status unix.WaitStatus
}

type GrimReaper struct {
	pidWaiters      map[int]PidWaitChannel
	wildcardWaiters []WildcardWaitChannel
}

func NewGrimReaper() *GrimReaper {
	reaper := new(GrimReaper)
	reaper.pidWaiters = map[int]PidWaitChannel{}
	reaper.wildcardWaiters = []WildcardWaitChannel{}

	return reaper
}

func (self *GrimReaper) WaitPid(pid int) PidWaitChannel {
	ch := make(PidWaitChannel)
	self.pidWaiters[pid] = ch

	return ch
}

func (self *GrimReaper) WaitWildcard() WildcardWaitChannel {
	ch := make(WildcardWaitChannel)
	self.wildcardWaiters = append(self.wildcardWaiters, ch)

	return ch
}

func (self *GrimReaper) Run() {
	signals := make(chan os.Signal)
	signal.Notify(signals, unix.SIGCHLD)

	for range signals {
		for {
			var status unix.WaitStatus
			pid, err := unix.Wait4(-1, &status, unix.WNOHANG, nil)

			if pid > 0 && err == nil {
				log.WithFields(log.Fields{
					"pid": pid,
				}).Debugf("Reaped process")

				if pidWaiter, exists := self.pidWaiters[pid]; exists {
					pidWaiter <- status
					delete(self.pidWaiters, pid)
				}

				for _, wildcardWaiter := range self.wildcardWaiters {
					wildcardWaiter <- WaitStatusPid{pid, status}
				}
			} else {
				break
			}
		}
	}
}
