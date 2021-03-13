# Verifyzero

This little utility ensures that a specified file is composed entirely of zero's. This application is intended to be used with magnetic block devices to ensure that they have been zeroized. It does not ensure that the device has been securely erased.

Usage: verifyzero /path/to/file

Exit Statuses
  0. Success, file is composed of zero's.
  1. Fail, non 0 bytes are contained within the file.
  2. Fail, an error has occoured. Details will be printed to stderr.

# Zeroizing a disk

My two favourite methods, note these are not secure methods of erasing disks they just zeroize. Also do not use either of these commands on an SSD.

1. dd if=/dev/zero of=/dev/sdz bs=4M
1. shred -n 0 -z -v /dev/sdz

Shred is preferred as it will provide a status indicator.