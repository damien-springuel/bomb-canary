package main

import (
	"log"
	"net"

	"github.com/damien-springuel/bomb-canary/server/generated"
	"github.com/damien-springuel/bomb-canary/server/lobby"
	"google.golang.org/grpc"
)

func main() {
	address := ":44324"
	grpcServer := grpc.NewServer()
	generated.RegisterLobbyServer(grpcServer, lobby.LobbyServer{})
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Error trying to listen to %s: %v", address, err)
	}
	log.Printf("Serving on %s", address)
	if err = grpcServer.Serve(listen); err != nil {
		log.Fatalf("Error serving grpc: %v", err)
	}
}
