package gapi

import (
	"context"
	"time"

	db "github.com/ShadrackAdwera/go-gRPC/db/sqlc"
	"github.com/ShadrackAdwera/go-gRPC/pb"
	"github.com/ShadrackAdwera/go-gRPC/workers"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ShadrackAdwera/go-gRPC/utils"
)

func (srv *Server) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {

	hashedPw, err := utils.HashPassword(req.GetPassword())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error hashing password: %s", err)
	}

	response, err := srv.store.CreateUserTx(ctx, db.CreateUserTxInput{
		CreateUserParams: db.CreateUserParams{
			Username: req.GetUsername(),
			Email:    req.GetEmail(),
			Password: hashedPw,
		},
		EmitCreateUser: func(user db.User) error {
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(workers.QueueCritical),
			}
			payload := workers.UserPayload{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			}
			return srv.distro.DistributeUser(ctx, payload, opts...)
		},
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "email exists: %s", pqErr)
			}
		}
		return nil, status.Errorf(codes.Internal, "error creating user %s", err)
	}

	return &pb.SignUpResponse{
		User: getUserResponse(response.User),
	}, nil
}

func getUserResponse(user db.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
