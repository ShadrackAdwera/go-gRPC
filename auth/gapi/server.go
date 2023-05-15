package gapi

import (
	db "github.com/ShadrackAdwera/go-gRPC/db/sqlc"
	"github.com/ShadrackAdwera/go-gRPC/pb"
	"github.com/ShadrackAdwera/go-gRPC/token"
	"github.com/ShadrackAdwera/go-gRPC/workers"
)

type Server struct {
	pb.UnimplementedAuthServer
	maker  token.TokenMaker
	store  db.TxStore
	distro workers.Distributor
}

func NewServer(maker token.TokenMaker, store db.TxStore, distro workers.Distributor) *Server {
	server := Server{
		maker:  maker,
		store:  store,
		distro: distro,
	}
	return &server
}
