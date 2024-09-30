//go:build linux || darwin
// +build linux darwin

package main

import (
	"log"
	"syscall"
)

func setMaxFiles(nofiles int) {
	var rLimit syscall.Rlimit

	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)

	if err != nil {
		log.Printf("Error Getting Rlimit %v", err)
	}

	rLimit.Max = uint64(nofiles)
	rLimit.Cur = uint64(nofiles)

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Printf("Error Setting Rlimit: %v", err)
		return
	}

	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Printf("Error Getting Rlimit %v", err)
		return
	}
	log.Printf("Successfully set max open file handles to %+v", rLimit)
}
