package main

import (
	"os"
	"testing"
)

func TestVerifyZero(t *testing.T) {
	file, err := os.Open("testdata/zero")
	if err != nil {
		t.Fatalf(err.Error())
	}

	status, _ := verifyZero(*file)

	if status == false {
		t.Fail()
	}
}

func TestVerifyNotZero(t *testing.T) {
	file, err := os.Open("testdata/notzero")
	if err != nil {
		t.Fatalf(err.Error())
	}

	status, _ := verifyZero(*file)

	if status == true {
		t.Fail()
	}
}
