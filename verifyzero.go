package main

import (
	"fmt"
	"io"
	"os"
)

func showhelp() {
	os.Stderr.WriteString("Usage: ./verifyzero /path/to/file\n")
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
				// Testing shows that upon error file.Read does not seek forward by the full len(buf)
				// On read error. I think it reads up to the error but I am not confident and this
				// suspicion will require furthur testing. Instead, after a read error I will seek
				// forward 1024 bytes and attempt the read again, while reporting each time that
				// an error has been encountered and the affected byte range.
				curpos, err := file.Seek(0, os.SEEK_CUR) // Get the current position
				if err != nil {
					os.Stderr.WriteString(err.Error())
					os.Exit(10)
				}

				newpos, err := file.Seek(1024, os.SEEK_CUR) // seek forward 1K
				if err != nil {
					os.Stderr.WriteString(err.Error())
					os.Exit(10)
				}

				os.Stderr.WriteString(fmt.Sprintf(": reading between offset %d-%d\n", curpos, newpos))

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
// 10 Fatal Error
func main() {
	if len(os.Args) < 2 {
		showhelp()
		os.Exit(10)
	}
	path := os.Args[1]
	exitstatus := 10

	file, err := os.Open(path)
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Exit(10)
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
