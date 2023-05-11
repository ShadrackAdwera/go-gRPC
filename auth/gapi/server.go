package gapi

import (
	db "github.com/ShadrackAdwera/go-gRPC/db/sqlc"
	"github.com/ShadrackAdwera/go-gRPC/pb"
	"github.com/ShadrackAdwera/go-gRPC/token"
)

type Server struct {
	pb.UnimplementedAuthServer
	maker token.TokenMaker
	store db.TxStore
}

func NewServer(maker token.TokenMaker, store db.TxStore) *Server {
	server := Server{
		maker: maker,
		store: store,
	}
	return &server
}
