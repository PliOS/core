// Copyright 2017 Saad Nasser (SdNssr). All rights reserved.
//
// Use of this source code is governed by a MIT-style license that can be found
// in the LICENSE file.

use std;
use nix;

/// A trait for handleable error types
pub trait HandleableError {
    /// If it is an error, print a message and loop forever
    fn handle_error(&self, message: &str);
}

impl<T, E> HandleableError for Result<T, E>
    where T: std::fmt::Debug,
          E: std::fmt::Debug
{
    fn handle_error(&self, message: &str) {
        if let &Err(ref e) = self {
            println!("Error - {}: {:?}", message, e);
            loop {}
        }
    }
}

/// Create a directory if it does not exist
pub fn create_directory(dir: &str, mode: u32) -> nix::Result<()> {
    use nix::unistd;
    use nix::sys::stat;
    use nix::Error::Sys;
    use nix::Errno::EEXIST;

    let error = unistd::mkdir(dir, stat::Mode::from_bits_truncate(mode));

    if error != Err(Sys(EEXIST)) {
        error
    } else {
        Ok(())
    }
}

/// Create a symlink if it does not exist
pub fn symlink_file(file: &str, to: &str) {
    use std::os::unix::fs;
    use std::io::ErrorKind::AlreadyExists;

    let result = fs::symlink(to, file);

    if let Err(err) = result {
        if err.kind() != AlreadyExists {
            println!("Error - Unable to symlink {} to {}: {:?}", file, to, err);
        }
    }
}
