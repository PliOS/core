# Core

PliOS core system

## init

A very simple init that only does the following:
- Sets up the environment ([environment.rs](init/src/environment.rs#L19))
- Mounts API filesystems - `/proc`, `sys`, etc. ([environment.rs](init/src/environment.rs#L65))
- Spawns programs based on signals ([event.rs](init/src/events.rs#L39))
