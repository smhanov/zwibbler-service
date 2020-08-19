package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const maxLogSize = 40 * 1024 * 1024

type logfile struct {
	file  *os.File
	name  string
	size  int64
	mutex sync.Mutex
}

func newLogFile(name string) *logfile {
	var file *os.File
	var err error

	lf := &logfile{}

	path := filepath.Dir(name)
	err = os.MkdirAll(path, 0776)
	if err != nil {
		log.Printf("Error: %v", err)
		return lf
	}

	_, err = os.Stat(name)

	if err != nil {
		file, err = os.Create(name)
	} else {
		file, err = os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}

	if err != nil {
		log.Printf("Error: %v", err)
		return lf
	}

	stat, err := file.Stat()
	if err != nil {
		log.Panic(err)
	}

	lf.file = file
	lf.size = stat.Size()
	lf.name = name
	log.Printf("Opened %s\n", lf.name)
	return lf
}

func (lf *logfile) Write(b []byte) (n int, err error) {
	lf.mutex.Lock()
	defer lf.mutex.Unlock()
	os.Stderr.Write(b)
	var written int
	if lf.file != nil {
		if lf.size > maxLogSize {
			// close file.
			lf.file.Close()

			// rename to old
			err := os.Rename(lf.name, lf.name+".old")
			if err != nil {
				log.Panic(err)
			}

			// open new file
			file, err := os.Create(lf.name)
			if err != nil {
				log.Panic(err)
			}

			lf.file = file
			lf.size = 0
		}
		lf.size += int64(len(b))
		written, err = lf.file.Write(b)
		if err != nil {
			fmt.Printf("Error writing to %s: %v\n", lf.name, err)
		}
	}

	return written, err
}
