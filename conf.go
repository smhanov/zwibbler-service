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
	compression bool

	// default: sqlite
	// can be redis or redis-cluster
	database string

	// default: localhost:6379
	redisServers  []string
	redisPassword string

	dbServers  []string
	dbUser     string
	dbPassword string

	secretUser     string
	secretPassword string

	jwtKey         string
	jwtKeyIsBase64 bool

	webhookURL string
	swarmUrls  []string

	maxFiles int64
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
	config.database = "sqlite"
	config.compression = true

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

		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}

		equals := strings.Index(line, "=")
		if equals >= 0 {
			key := strings.TrimSpace(line[:equals])
			value := strings.TrimSpace(line[equals+1:])
			switch key {
			case "ServerPort":
				i, _ := strconv.ParseInt(value, 10, 32)
				config.port = int(i)

			case "ServerBindAddress":
				config.bindAddress = value

			case "CertFile":
				if value != "" && runtime.GOOS == "windows" && fileExists("\\zwibbler\\"+value) {
					value = "\\zwibbler\\" + value
				}
				config.certFile = value
			case "KeyFile":
				if value != "" && runtime.GOOS == "windows" && fileExists("\\zwibbler\\"+value) {
					value = "\\zwibbler\\" + value
				}
				config.keyFile = value
			case "Expiration":
				value = strings.ToLower(value)
				if value == "never" {
					config.expiration = zwibserve.NoExpiration
				} else {
					config.expiration, _ = strconv.ParseInt(value, 10, 64)
				}
			case "Database":
				switch value {
				case "redis", "redis-cluster", "sqlite", "postgres", "mysql", "mariadb":
					config.database = value
				default:
					log.Printf("Error: Unknown database type %s, must be redis, sqlite, postgres, mysql, or mariadb", value)
				}
			case "RedisServer":
				config.dbServers = append(config.dbServers, value)
			case "RedisPassword":
				config.dbPassword = value
			case "DbServer":
				config.dbServers = append(config.dbServers, value)
			case "DbUser":
				config.dbUser = value
			case "DbPassword":
				config.dbPassword = value
			case "Compression":
				value = strings.ToLower(value)
				config.compression = isTrue(value)
			case "SecretUser":
				config.secretUser = value
			case "SecretPassword":
				config.secretPassword = value
			case "JWTKey":
				config.jwtKey = value
			case "JWTKeyIsBase64":
				config.jwtKeyIsBase64 = isTrue(value)
			case "Webhook":
				config.webhookURL = value
			case "MaxFiles":
				config.maxFiles, _ = strconv.ParseInt(value, 10, 64)
			case "SwarmURL":
				if len(value) > 0 {
					config.swarmUrls = append(config.swarmUrls, value)
				}
			}
		}
	}

	if len(config.redisServers) == 0 {
		switch config.database {
		case "redis":
			config.dbServers = []string{"localhost:6379"}
		case "mysql":
			config.dbServers = []string{"localhost:3306"}
		case "postgres":
			config.dbServers = []string{"localhost:5432"}
		case "mariadb":
			config.dbServers = []string{"localhost:3306"}
		}
	}

	return config, nil
}

func isTrue(value string) bool {
	value = strings.ToLower(value)
	return value != "0" && value != "false" && value != "off"
}
