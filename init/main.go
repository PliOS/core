package main

import (
	"golang.org/x/sys/unix"
	"log"
	"os"
	"os/exec"
	"os/signal"
)

func SpawnProcess(name string, arg ...string) int {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("Fatal error - running %s: %s", name, err)
	}

	return cmd.Process.Pid
}

func SpawnInitEvent() {
	// This is run in a goroutine as we want to wait a bit before we start it
	// It does nothing right now, but it should
	// time.Sleep(100 * time.Millisecond)
	// SpawnProcess("systemctl", "init")
}

func ReapChildren() {
	serviceManagerPid := SpawnProcess("sh")

	go SpawnInitEvent()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, unix.SIGCHLD)

	for {
		select {
		case <-signals:
		}

		for {
			var status unix.WaitStatus
			pid, err := unix.Wait4(-1, &status, unix.WNOHANG, nil)

			switch err {
			case nil:
				if pid == serviceManagerPid {
					serviceManagerPid = SpawnProcess("sh")
				}

				if pid <= 0 {
					break
				}
			default:
				break
			}
		}
	}
}

func main() {
	log.SetPrefix("[init] ")

	SetupProcessEnvironment()
	MountApiFilesystems()
	CreateApiSymlinks()

	log.Printf("Finished system setup\n")

	signals := make(chan os.Signal, 1)
	NotifyAllSignals(signals)

	go ReapChildren()

	for signal := range signals {
		switch signal {
		case PEACEFUL_HALT:
			SpawnProcess("systemctl", "halt")
		case PEACEFUL_SHUTDOWN:
			SpawnProcess("systemctl", "shutdown")
		case PEACEFUL_REBOOT:
			SpawnProcess("systemctl", "reboot")
		case FORCEFUL_HALT:
			unix.Reboot(unix.LINUX_REBOOT_CMD_HALT)
		case FORCEFUL_SHUTDOWN:
			unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF)
		case FORCEFUL_REBOOT:
			unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART)
		case RECOVERY:
			SpawnProcess("systemctl", "recovery")
		case RESCUE:
			SpawnProcess("systemctl", "rescue")
		case NORMAL:
			SpawnProcess("systemctl", "normal")
		case POWERFAIL:
			SpawnProcess("systemctl", "powerfail")
		case KBREQ:
			SpawnProcess("systemctl", "kbreq")
		case SAK:
			SpawnProcess("systemctl", "saq")
		}
	}
}
