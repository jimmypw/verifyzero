Verifyzero

This little utility ensures that a specified file is composed entirely of zero's. This application is intended to be used with magnetic block devices to ensure that they have been zeroized. It does not ensure that the device has been securely erased.

Usage: verifyzero /path/to/file

Exit Status'
0. Success, file is composed of errors
1. Fail, non 0 bytes are contained within the file
2. Fail, an error has occoured. Details will be printed to stderr.