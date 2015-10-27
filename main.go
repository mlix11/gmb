package main

import (
	"flag"
	"github.com/loose11/gmb/config"
	"github.com/loose11/gmb/database"
	"github.com/loose11/gmb/handler"
	"net/http"
)

func init() {
	flag.StringVar(&config.Port, "port", config.DefaultPort, "Default Port")
	flag.StringVar(&config.BasePath, "basePath", config.DefaultBasePath, "Default BasePath")
}

func main() {
	flag.Parse()

	db := database.NewAppDatabase(config.BasePath)

	handler := handler.HandlerConfig{db}

	http.HandleFunc("/", handler.Standard)
	http.ListenAndServe(":"+config.Port, nil)
}
