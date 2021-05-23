package main

import (
	"flag"
	"log"
	"russian-learns-english/internal/api/component/yandex"
	"russian-learns-english/internal/api/delivery/http/server"
	"russian-learns-english/internal/api/service"
	"strconv"
)

var (
	port uint
)

func init() {
	flag.UintVar(&port, "port", 8080, "Port that HTTP server will start on")
	flag.StringVar(&configPath, "config-path", "api.yml", "Path to config file")
	flag.Parse()
}

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}

	yaDict := &yandex.DictAPIClient{ApiKey: cfg.YandexDictApiKey}
	if _, err = yaDict.TranslateWord("text"); err != nil {
		log.Fatalf("Invalid yandex token")
	}

	srv := server.New(service.NewWordService(yaDict))
	srv.Run(":" + strconv.Itoa(int(port)))
}
