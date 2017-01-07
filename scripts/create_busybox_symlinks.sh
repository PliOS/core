#!/bin/sh

# Copyright 2017 Saad Nasser (SdNssr). All rights reserved.
#
# Use of this source code is governed by a MIT-style license that can be found
# in the LICENSE file.

applets="acpid add-shell addgroup adduser adjtimex arp arping ash awk base64 basename beep    \
blkdiscard blkid blockdev bootchartd brctl bunzip2 bzcat bzip2 cal cat catv chat     \
chattr chgrp chmod chown chpasswd chpst chroot chrt chvt cksum clear cmp comm conspy \
cp cpio crond crontab cryptpw cttyhack cut date dc dd deallocvt delgroup deluser     \
depmod devmem df dhcprelay diff dirname dmesg dnsd dnsdomainname dos2unix dpkg       \
dpkg-deb du dumpkmap dumpleases echo ed egrep eject env envdir envuidgid ether-wake  \
expand expr fakeidentd false fatattr fbset fbsplash fdflush fdformat fdisk fgconsole \
fgrep find findfs flock fold free freeramdisk fsck fsck.minix fstrim fsync ftpd      \
ftpget ftpput fuser getopt getty grep groups gunzip gzip halt hd hdparm head hexdump \
hostid hostname httpd hush hwclock i2cdetect i2cdump i2cget i2cset id ifconfig       \
ifdown ifenslave ifplugd ifup inetd insmod install ionice iostat ip ipaddr ipcalc    \
ipcrm ipcs iplink ipneigh iproute iprule iptunnel kbd_mode kill killall killall5     \
klogd less linux32 linux64 linuxrc ln loadfont loadkmap logger login logname logread \
losetup lpd lpq lpr ls lsattr lsmod lsof lspci lsusb lzcat lzma lzop lzopcat         \
makedevs makemime man md5sum mdev mesg microcom mkdir mkdosfs mke2fs mkfifo          \
mkfs.ext2 mkfs.minix mkfs.vfat mknod mkpasswd mkswap mktemp modinfo modprobe more    \
mount mountpoint mpstat mt mv nameif nanddump nandwrite nbd-client nc netstat nice   \
nmeter nohup nsenter nslookup ntpd od openvt passwd patch pgrep pidof ping ping6     \
pipe_progress pivot_root pkill pmap popmaildir poweroff powertop printenv printf ps  \
pscan pstree pwd pwdx raidautorun rdate rdev readahead readlink readprofile realpath \
reboot reformime remove-shell renice reset resize rev rm rmdir rmmod route rpm       \
rpm2cpio rtcwake run-parts runsv runsvdir rx script scriptreplay sed sendmail seq    \
setarch setconsole setfont setkeycodes setlogcons setserial setsid setuidgid sh      \
sha1sum sha256sum sha3sum sha512sum showkey shuf slattach sleep smemcap softlimit    \
sort split start-stop-daemon stat strings stty su sulogin sum sv svc svlogd swapoff  \
swapon switch_root sync sysctl syslogd tac tail tar tcpsvd tee telnet telnetd test   \
tftp tftpd time timeout top touch tr traceroute traceroute6 true truncate tty        \
ttysize tunctl ubiattach ubidetach ubimkvol ubirename ubirmvol ubirsvol ubiupdatevol \
udhcpc udhcpd udpsvd uevent umount uname unexpand uniq unix2dos unlink unlzma unlzop \
unshare unxz unzip uptime usleep uudecode uuencode vconfig vi vlock volname watch    \
watchdog wc wget which whoami whois xargs xz xzcat yes zcat zcip"


for applet in $applets
do
      sudo ln /media/plios_sysroot/bin/busybox /media/plios_sysroot/bin/$applet
done
