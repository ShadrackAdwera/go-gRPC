package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	db "github.com/ShadrackAdwera/go-gRPC/db/sqlc"
	"github.com/ShadrackAdwera/go-gRPC/gapi"
	"github.com/ShadrackAdwera/go-gRPC/pb"
	"github.com/ShadrackAdwera/go-gRPC/token"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}
	url := os.Getenv("PG_URL")
	conn, err := sql.Open("postgres", url)

	if err != nil {
		log.Fatalf("Failed to initialize the database %v", err)
	}
	symmKey := os.Getenv("SYMMETRIC_KEY")

	maker, err := token.NewPasetoMaker(symmKey)

	if err != nil {
		log.Fatalf("Failed to create the token maker %v", err)
	}

	store := db.NewStore(conn)

	server := gapi.NewServer(maker, store)

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer, server)
	reflection.Register(grpcServer) // allows client to inspect available RPCs on the server and how to call them
	addr := os.Getenv("GRPC_SERVER_ADDRESS")

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Failed to create the grpc listener %v", err)
	}

	log.Printf("start gRPC server on PORT: %s\n", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Failed to serve gRPC requests %v", err)
	}

}
