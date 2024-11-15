package main

import (
	"flag"
	"strings"
)

type config struct {
	isProd             bool
	port               int
	allowedOrigins     []string
	frontendBundlePath string
}

func GetConfig() config {
	isProd := flag.Bool("prod", false, "set the server to production mode")
	port := flag.Int("port", 44324, "server port")
	allowedOriginsString := flag.String("allowed-origins", "http://localhost:44322", "define the allowed origins for CORS")
	allowedOrigins := strings.Split(*allowedOriginsString, ",")
	frontendBundlePath := flag.String("fe-bundle-path", "../client/dist", "set the path to the frontend bundle")

	flag.Parse()

	return config{
		isProd:             *isProd,
		port:               *port,
		allowedOrigins:     allowedOrigins,
		frontendBundlePath: *frontendBundlePath,
	}
}
