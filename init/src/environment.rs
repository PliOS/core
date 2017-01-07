// Copyright 2017 Saad Nasser (SdNssr). All rights reserved.
//
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

use nix;
use libc;

use std::env;

use nix::mount;
use nix::sys::stat;

use utils;
use utils::HandleableError;

/// Initialize the bare minimum linux environment
///
/// - Disables ctrl-alt-delete reboot
/// - setsid()
/// - chdir('/')
/// - umask(0000)
/// - Sets PATH and LANG
pub fn init_environment() {
    let error = nix::sys::reboot::set_cad_enabled(false);

    error.handle_error("Unable to disable ctrl-alt-delete reboot");

    unsafe {
        libc::setsid();
    }

    let error = env::set_current_dir("/");

    error.handle_error("Unable to change directory to /");

    stat::umask(stat::Mode::empty());

    env::set_var("PATH",
                 "/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin");
    env::set_var("RUST_BACKTRACE", "1");
    env::set_var("LANG", "C");
}

/// Mounts a filesystem in a directory with permissions
///
/// Creates the directory if it does not exist.
fn init_api_filesystem(source: &str,
                       mountpoint: &str,
                       fstype: &str,
                       flags: mount::MsFlags,
                       data: Option<&str>,
                       permissions: u32)
                       -> nix::Result<()> {
    // Create the mount directory if it does not exist
    try!(utils::create_directory(mountpoint, permissions));
    mount::mount(Some(source), mountpoint, Some(fstype), flags, data)
}

/// Mount all API filesystems
///
/// - /proc
/// - /sys
/// - /sys/fs/cgroup
/// - /sys/fs/cgroup/systemd
/// - /run
/// - /run/shm
/// - /dev
/// - /dev/pts
pub fn init_api_filesystems() {
    println!("Mounting /proc...");

    let error = init_api_filesystem("proc",
                                    "/proc",
                                    "proc",
                                    mount::MS_NOSUID | mount::MS_NODEV,
                                    None,
                                    0o555);

    error.handle_error("Unable to mount /proc");

    println!("Mounting /sys...");

    let error = init_api_filesystem("sysfs",
                                    "/sys",
                                    "sysfs",
                                    mount::MS_NOSUID | mount::MS_NODEV | mount::MS_NOEXEC,
                                    None,
                                    0o555);

    error.handle_error("Unable to mount /sys");

    println!("Mounting /sys/fs/cgroup...");

    let error = init_api_filesystem("tmpfs",
                                    "/sys/fs/cgroup",
                                    "tmpfs",
                                    mount::MS_NOSUID | mount::MS_NODEV | mount::MS_NOEXEC,
                                    Some("size=1M,mode=0755"),
                                    0o0755);

    error.handle_error("Unable to mount /sys/fs/cgroup");

    println!("Mounting /sys/fs/cgroup/systemd...");

    let error = init_api_filesystem("cgroup",
                                    "/sys/fs/cgroup/systemd",
                                    "cgroup",
                                    mount::MS_NOSUID | mount::MS_NODEV | mount::MS_NOEXEC,
                                    Some("name=systemd,none"),
                                    0o0755);

    error.handle_error("Unable to mount /sys/fs/cgroup/systemd");

    println!("Mounting /run...");

    let error = init_api_filesystem("tmpfs",
                                    "/run",
                                    "tmpfs",
                                    mount::MS_NOSUID | mount::MS_NODEV | mount::MS_STRICTATIME,
                                    Some("size=20%,mode=0755"),
                                    0o0755);

    error.handle_error("Unable to mount /run");

    println!("Mounting /run/shm...");

    let error = init_api_filesystem("tmpfs",
                                    "/run/shm",
                                    "tmpfs",
                                    mount::MS_NOSUID | mount::MS_NODEV | mount::MS_NOEXEC |
                                    mount::MS_STRICTATIME,
                                    Some("size=50%,mode=01777"),
                                    0o01777);

    error.handle_error("Unable to mount /run/shm");

    println!("Mounting /dev...");

    let error = init_api_filesystem("devtmpfs",
                                    "/dev",
                                    "devtmpfs",
                                    mount::MS_NOSUID | mount::MS_STRICTATIME,
                                    Some("size=10M,mode=0755"),
                                    0o0755);

    error.handle_error("Unable to mount /dev");

    println!("Mounting /dev/pts...");

    let error = init_api_filesystem("devpts",
                                    "/dev/pts",
                                    "devpts",
                                    mount::MS_NOSUID | mount::MS_NOEXEC,
                                    Some("ptmxmode=0666,gid=5,newinstance,mode=0620"),
                                    0o0620);

    error.handle_error("Unable to mount /dev/pts");

    utils::symlink_file("/dev/ptmx", "/dev/pts/ptmx");
    utils::symlink_file("/dev/fd", "/proc/self/fd");
    utils::symlink_file("/dev/core", "/proc/kcore");
    utils::symlink_file("/dev/stdin", "/proc/self/fd/0");
    utils::symlink_file("/dev/stdout", "/proc/self/fd/1");
    utils::symlink_file("/dev/stderr", "/proc/self/fd/2");
    utils::symlink_file("/dev/shm", "/run/shm");
}
