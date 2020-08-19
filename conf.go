package main

import (
	"bufio"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type configFile struct {
	bindAddress string
	port        int
	certFile    string
	keyFile     string
}

func readConfFile() (configFile, error) {
	var config configFile
	config.bindAddress = "0.0.0.0"
	config.port = 3000

	confPath := "/etc/zwibbler.conf"
	if runtime.GOOS == "windows" {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		confPath = path.Join(dir, "zwibbler.conf")
	}

	// read file,
	file, err := os.Open(confPath)
	if err != nil {
		log.Printf("Could not open conf file %s: %v", confPath, err)
		return config, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "=")
		if strings.HasPrefix(line[0], "#") {
			continue
		}

		if len(line) == 2 {
			key := strings.TrimSpace(line[0])
			value := strings.TrimSpace(line[1])
			if key == "ServerPort" {
				i, _ := strconv.ParseInt(value, 10, 32)
				config.port = int(i)
			} else if key == "ServerBindAddress" {
				config.bindAddress = value
			} else if key == "CertFile" {
				config.certFile = value
			} else if key == "KeyFile" {
				config.keyFile = value
			}
		}
	}

	return config, nil
}
