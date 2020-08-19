package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/kardianos/service"
	"github.com/smhanov/zwibserve"
)

type program struct {
	server *http.Server
}

func (p *program) run() {
	config, err := readConfFile()
	if err != nil {
		log.Panic(err)
	}
	http.Handle("/socket", zwibserve.NewHandler(zwibserve.NewSQLITEDB("/var/lib/zwibbler.db")))
	bind := fmt.Sprintf("%s:%d", config.bindAddress, config.port)

	p.server = &http.Server{
		Addr: bind,
	}

	if config.certFile == "" {
		log.Printf("Listening on %s...", bind)
		err = p.server.ListenAndServe()
	} else {
		// HTTPS server!
		log.Printf("Listening on HTTPS %s...", bind)
		p.server.ListenAndServeTLS(config.certFile, config.keyFile)
	}

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	p.server.Shutdown(context.Background())
	return nil
}
