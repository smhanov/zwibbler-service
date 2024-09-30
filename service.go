package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/go-redis/redis/v9"
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

	if overridePort != 0 {
		config.port = overridePort
	}

	var dbpath = "/var/lib/zwibbler/"
	if runtime.GOOS == "windows" {
		dbpath = "\\zwibbler\\"
	}

	err = os.MkdirAll(dbpath, 0776)
	if err != nil {
		log.Panic(err)
	}

	if config.maxFiles > 0 {
		setMaxFiles(int(config.maxFiles))
	}

	var db zwibserve.DocumentDB

	switch config.database {
	case "sqlite":
		log.Printf("Database path is %s", dbpath+"zwibbler.db")
		db = zwibserve.NewSQLITEDB(dbpath + "zwibbler.db")
	case "postgres":
		log.Printf("Using PostgresSQL database %s", config.dbServers[0])
		db = zwibserve.NewPostgreSQLConnection(config.dbServers[0], config.dbUser, config.dbPassword, "zwibbler")
	case "mysql":
		log.Printf("Using MySQL database %s", config.dbServers[0])
		db = zwibserve.NewMySQLConnection(config.dbServers[0], config.dbUser, config.dbPassword, "zwibbler")
	case "mariadb":
		log.Printf("Using MariaDB database %s", config.dbServers[0])
		db = zwibserve.NewMariaDBConnection(config.dbServers[0], config.dbUser, config.dbPassword, "zwibbler")
	case "redis":
		log.Printf("Using Redis DB %s", config.dbServers[0])
		db = zwibserve.NewRedisDB(&redis.Options{
			Addr:     config.dbServers[0],
			Password: config.dbPassword,
		})
	case "redis-cluster":
		log.Printf("Using Redis Cluster DB %v", config.dbServers)
		db = zwibserve.NewRedisClusterDB(&redis.ClusterOptions{
			Addrs:    config.dbServers,
			Password: config.dbPassword,
		})
	}

	if lostWriteRate > 0 {
		log.Printf("Simulating unreliable database with 1 in %d chance of silently dropping the write", lostWriteRate)
		db = createUnreliableDB(db, lostWriteRate)
	}

	if config.expiration == 0 {
		log.Printf("Set document expiration to 24 hours (default)")
	} else if config.expiration == zwibserve.NoExpiration {
		log.Printf("Set expiration to NEVER")
	} else {
		log.Printf("Set expiration to %d seconds", config.expiration)
	}
	db.SetExpiration(config.expiration)

	handler := zwibserve.NewHandler(db)

	log.Printf("Socket compression allowed: %v", config.compression)
	handler.SetCompressionAllowed(config.compression)

	if config.secretUser != "" || config.secretPassword != "" {
		log.Printf("Secret user specified. Managment API will be enabled.")
		handler.SetSecretUser(config.secretUser, config.secretPassword)
	}

	if config.jwtKey != "" {
		log.Printf("JWTKey specified. Only JWT will be accepted from clients.")
		handler.SetJWTKey(config.jwtKey, config.jwtKeyIsBase64)
	}

	handler.SetSwarmURLs(config.swarmUrls)

	log.Printf("Webhook URL: %v", config.webhookURL)
	handler.SetWebhookURL(config.webhookURL)

	http.Handle("/socket", handler)
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
