package main

import (
	"flag"
	"strconv"
)

var IsProd string

type config struct {
	isProd             bool
	port               int
	allowedOrigins     []string
	frontendBundlePath string
}

func GetConfig() config {
	isProd, err := strconv.ParseBool(IsProd)
	if err != nil {
		isProd = false
	}

	portFlag := flag.Int("port", 44333, "server port")
	flag.Parse()
	port := *portFlag

	frontendBundlePath := "."
	allowedOrigins := []string{"*"}
	if !isProd {
		frontendBundlePath = "../client/dist"
		port = 44324
		allowedOrigins = []string{"http://localhost:44322"}
	}

	return config{
		isProd:             isProd,
		port:               port,
		allowedOrigins:     allowedOrigins,
		frontendBundlePath: frontendBundlePath,
	}
}
