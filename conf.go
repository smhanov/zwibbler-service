package main

import (
	"bufio"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/smhanov/zwibserve"
)

type configFile struct {
	bindAddress string
	port        int
	certFile    string
	keyFile     string
	expiration  int64
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func readConfFile() (configFile, error) {
	var config configFile
	config.bindAddress = "0.0.0.0"
	config.port = 3000

	confPath := "/etc/zwibbler.conf"
	if runtime.GOOS == "windows" {
		confPath = "\\zwibbler\\zwibbler.conf"
	}

	log.Printf("Reading configuration file %s", confPath)

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
				if value != "" && runtime.GOOS == "windows" && fileExists("\\zwibbler\\"+value) {
					value = "\\zwibbler\\" + value
				}
				config.certFile = value
			} else if key == "KeyFile" {
				if value != "" && runtime.GOOS == "windows" && fileExists("\\zwibbler\\"+value) {
					value = "\\zwibbler\\" + value
				}
				config.keyFile = value
			} else if key == "Expiration" {
				value = strings.ToLower(value)
				if value == "never" {
					config.expiration = zwibserve.NoExpiration
				} else {
					config.expiration, _ = strconv.ParseInt(value, 10, 64)
				}
			}
		}
	}

	return config, nil
}
