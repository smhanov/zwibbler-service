package main

import (
	"flag"
	"log"
	"strings"

	"github.com/kardianos/service"
)

func main() {
	svcConfig := &service.Config{
		Name:        "zwibbler",
		DisplayName: "Zwibbler Collaboration Service",
		Description: "Provides collaboration features for Zwibbler",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	var install bool
	var uninstall bool

	flag.BoolVar(&install, "install", false, "Install")
	flag.BoolVar(&uninstall, "uninstall", false, "Uninstall")
	flag.Parse()

	if install {
		log.Printf("Installing service.")
		err = s.Install()
		if err != nil && strings.Index(err.Error(), "already exists") < 0 {
			log.Panic(err)
		}

		log.Printf("Restarting zwibbler service")
		if status, _ := s.Status(); status == service.StatusRunning {
			err = s.Restart()
		} else {
			log.Printf("Starting zwibbler service")
			err = s.Start()
		}
	}

	if uninstall {
		log.Printf("Uninstalling service.")
		err = s.Uninstall()
		if err != nil {
			log.Printf("Error uninstalling: %v", err)
		}
		return
	}

	if err != nil {
		log.Panic(err)
	}

	if !install && !uninstall {
		serverLog := "/var/log/zwibbler/zwibbler.log"
		logger := newLogFile(serverLog)
		log.SetOutput(logger)

		err = s.Run()
		if err != nil {
			log.Panic(err)
		}
	}

}
