package main

import (
	"io"
	"os"
)

func showhelp() {
	os.Stderr.WriteString("Usage: verifyzero /path/to/file\n")
}

func verifyZero(file os.File) (bool, error) {
	var buf = make([]byte, 1024*1000*10) // 10MB

	for {
		readlen, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return false, err
			}
		}

		for i := 0; i < readlen; i++ {
			if buf[i] != 0x00 {
				return false, nil
			}
		}

	}
	return true, nil
}

// Exit status'
// 0 success (target is empty)
// 1 fail (target is not empty)
// 2 fail (there was an error)
func main() {
	if len(os.Args) < 2 {
		showhelp()
		os.Exit(2)
	}
	path := os.Args[1]
	exitstatus := 1

	file, err := os.Open(path)
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Exit(2)
	}

	result, err := verifyZero(*file)
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Exit(2)
	}

	file.Close()

	if result {
		exitstatus = 0
	} else {
		exitstatus = 1
	}

	os.Exit(exitstatus)
}
