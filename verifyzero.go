package main

import (
	"fmt"
	"io"
	"os"
)

func showhelp() {
	os.Stderr.WriteString("Usage: verifyzero /path/to/file\n")
}

func verifyZero(file os.File) (bool, bool) {
	var buf = make([]byte, 1024*1000*10) // 10MB
	fileiszeroized := true
	haserrors := false

	for {
		readlen, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				haserrors = true
				os.Stderr.WriteString(err.Error())
				// Testing shows that upon error file.Read does not seek forward. I therefore need to
				// step through the file until i identify the exact block that is causing the problem
				// then step over it. However this algorithm will be ineffienent so until I can think
				// of a way to do this properly I'll Seek foward 10MB and report the range where the
				// Error Occoured
				curpos, err := file.Seek(0, os.SEEK_CUR)
				if err != nil {
					os.Stderr.WriteString(err.Error())
					os.Exit(2)
				}

				newpos, err := file.Seek(1024*1000*10, os.SEEK_CUR) // 10MB
				if err != nil {
					os.Stderr.WriteString(err.Error())
					os.Exit(2)
				}

				os.Stderr.WriteString(fmt.Sprintf("Error reading between offset %d-%d\n", curpos, newpos))

				continue
			}
		}

		for i := 0; i < readlen; i++ {
			if buf[i] != 0x00 {
				fileiszeroized = false
				return fileiszeroized, haserrors
			}
		}

	}
	return fileiszeroized, haserrors
}

// Exit status'
// 0 success (target is empty)
// 1 success (with read errors)
// 2 fail (target is not empty)
// 3 fail (target is not empty and there was an error)
func main() {
	if len(os.Args) < 2 {
		showhelp()
		os.Exit(2)
	}
	path := os.Args[1]
	exitstatus := 2

	file, err := os.Open(path)
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Exit(2)
	}

	iszeroized, haserrors := verifyZero(*file)

	file.Close()

	if iszeroized && !haserrors {
		exitstatus = 0
	} else if iszeroized && haserrors {
		exitstatus = 1
	} else if !iszeroized && !haserrors {
		exitstatus = 2
	} else if !iszeroized && haserrors {
		exitstatus = 3
	}
	os.Exit(exitstatus)
}
