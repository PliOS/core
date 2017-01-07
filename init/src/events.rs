use std;
use std::mem;
use std::process::Command;

use libc;

use nix::sys::signal::SigSet;

use utils::HandleableError;

/// It's hardcoded right now. For whatever reason, libc does not have any mention
/// of realtime signals
const SIGRTMIN: i32 = 34;

const PEACEFUL_HALT: i32 = SIGRTMIN + 3;
const PEACEFUL_SHUTDOWN: i32 = SIGRTMIN + 4;
const PEACEFUL_REBOOT: i32 = SIGRTMIN + 5;

const FORCEFUL_HALT: i32 = SIGRTMIN + 13;
const FORCEFUL_SHUTDOWN: i32 = SIGRTMIN + 14;
const FORCEFUL_REBOOT: i32 = SIGRTMIN + 15;

const RECOVERY: i32 = SIGRTMIN + 2;
const RESCUE: i32 = SIGRTMIN + 1;
const NORMAL: i32 = SIGRTMIN + 0;

const POWERFAIL: i32 = libc::SIGPWR;
const KBREQ: i32 = libc::SIGWINCH;
const SAK: i32 = libc::SIGINT;

const REAP: i32 = libc::SIGCHLD;

/// Loop forever processing signals.
pub fn handle_events() {
    let mut service_manager_pid: i32;
    let signals = SigSet::all();

    signals.thread_set_mask().handle_error("Unable to set signal mask");

    // Spawn the service manager
    service_manager_pid = spawn_service_manager() as i32;

    // First event
    println!("systemctl target init");

    loop {
        match wait_signal() {
            PEACEFUL_HALT => {
                println!("systemctl target halt");
            }
            PEACEFUL_SHUTDOWN => {
                println!("systemctl target shutdown");
            }
            PEACEFUL_REBOOT => {
                println!("systemctl target reboot");
            }
            FORCEFUL_HALT => {
                println!("forceful halt");
            }
            FORCEFUL_SHUTDOWN => {
                println!("forceful shutdown");
            }
            FORCEFUL_REBOOT => {
                println!("forceful reboot");
            }
            RECOVERY => {
                println!("systemctl target recovery");
            }
            RESCUE => {
                println!("systemctl target rescue");
            }
            NORMAL => {
                println!("systemctl target normal");
            }
            POWERFAIL => {
                println!("systemctl target powerfail");
            }
            KBREQ => {
                println!("systemctl target kbreq");
            }
            SAK => {
                println!("systemctl target saq");
            }
            REAP => {
                if reap_child(service_manager_pid) {
                    println!("Service manager died...");
                    println!("Respawning service manager...");
                    service_manager_pid = spawn_service_manager() as i32;
                }
            }
            _ => {
                println!("Received unknown signal");
            }
        }
    }
}

fn spawn_service_manager() -> u32 {
    let child = spawn_process("sh", &[]);
    child.handle_error("Failed to spawn service_manager");
    return child.unwrap();
}

fn spawn_process(program: &str, arguments: &[&str]) -> Result<u32, std::io::Error> {
    Ok(try!(Command::new(program).args(arguments).spawn()).id())
}

fn reap_child(service_manager_pid: i32) -> bool {
    let mut did_service_manager_die = false;

    let mut status: i32 = unsafe { mem::uninitialized() };
    let mut pid: i32 = 1;

    while pid > 0 {
        pid = unsafe { libc::waitpid(-1, &mut status as *mut i32, libc::WNOHANG) };

        if pid == service_manager_pid {
            did_service_manager_die = true;
        }
    }

    did_service_manager_die
}

fn wait_signal() -> i32 {
    let mut sigset: libc::sigset_t = unsafe { mem::uninitialized() };
    let _ = unsafe { libc::sigfillset(&mut sigset as *mut libc::sigset_t) };

    let mut signum: libc::c_int = unsafe { mem::uninitialized() };
    let _ = unsafe { libc::sigwait(&sigset as *const libc::sigset_t, &mut signum) };

    return signum as i32;
}
