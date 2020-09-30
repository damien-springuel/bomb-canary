package main

import (
	"log"
	"net/http"

	"github.com/damien-springuel/bomb-canary/server/lobby"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	lobby.Register(router)

	port := ":44324"
	log.Printf("serving %s\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("error serving %v\n", err)
	}
}
