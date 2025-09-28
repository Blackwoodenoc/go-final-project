package server

import (
	"net/http"
 	"os"
 	"go1f/pkg/api"
)


func Run() error {
	port := "7540"
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}

	webDir := "web"
	fileServer := http.FileServer(http.Dir(webDir))

	api.Init() // регистрируем API

	http.Handle("/", fileServer) // статика
	
	return http.ListenAndServe(":"+port, nil)
}