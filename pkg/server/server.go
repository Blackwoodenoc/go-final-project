package server

import (
    "net/http"
    "os"
)

func Run() error{
	port := "7540"
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
 	}

 	webDir := "web"

	err := http.ListenAndServe(":"+port, http.FileServer(http.Dir(webDir)))
 	if err != nil {
  		panic(err)
 	}
	return nil
}