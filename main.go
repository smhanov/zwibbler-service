package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/kardianos/service"
	"github.com/smhanov/zwibserve"
)

func main() {
	var install bool
	var uninstall bool
	var test string
	var teachers int
	var students int
	var doc string
	var verbose bool

	flag.BoolVar(&install, "install", false, "Install")
	flag.BoolVar(&uninstall, "uninstall", false, "Uninstall")
	flag.StringVar(&test, "test", "", "Give wss:// url of another server to stress test")
	flag.StringVar(&doc, "docid", "", "Document ID for stress test")
	flag.IntVar(&teachers, "teachers", 0, "Number of teachers for stress test")
	flag.IntVar(&students, "students", 0, "Number of students for stress test")
	flag.BoolVar(&verbose, "verbose", false, "Verbose server")
	flag.Parse()

	if test != "" {
		if doc == "" {
			fmt.Printf("Error: You must specify the document id using --docid")
			return
		}
		if teachers == 0 && students == 0 {
			fmt.Printf("Error: You must specify non-zero --teachers or --students")
			return
		}

		// run a stress test
		zwibserve.RunStressTest(zwibserve.StressTestArgs{
			Address:     test,
			NumTeachers: teachers,
			NumStudents: students,
			DocumentID:  doc,
			Verbose:     verbose,
		})

		return
	}

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

	if install {
		log.Printf("Installing service.")
		err = s.Install()
		if err != nil && !strings.Contains(err.Error(), "already exists") {
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
		if status, _ := s.Status(); status == service.StatusRunning {
			log.Printf("Stopping service.")
			err = s.Stop()
			if err != nil {
				log.Printf("Error stopping: %v", err)
			}
		}

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
		if runtime.GOOS == "windows" {
			serverLog = "\\zwibbler\\zwibbler.log"
			os.MkdirAll("\\zwibbler", 0776)
		}

		logger := newLogFile(serverLog)
		log.SetOutput(logger)

		err = s.Run()
		if err != nil {
			log.Panic(err)
		}
	}

}
