package gapi

import (
	"testing"

	db "github.com/ShadrackAdwera/go-gRPC/db/sqlc"
	"github.com/ShadrackAdwera/go-gRPC/token"
	"github.com/ShadrackAdwera/go-gRPC/workers"
	"github.com/stretchr/testify/require"
)

const symmKey = "1185489AE92431DBA8E4C4BC2EA55241"

func newTestServer(t *testing.T, store db.TxStore, distro workers.Distributor) *Server {
	maker, err := token.NewPasetoMaker(symmKey)

	require.NoError(t, err)

	srv := NewServer(maker, store, distro)
	return srv
}
