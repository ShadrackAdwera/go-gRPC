package gapi

import (
	"context"

	db "github.com/ShadrackAdwera/go-gRPC/db/sqlc"
	"github.com/ShadrackAdwera/go-gRPC/pb"
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

	// check if user exists
	// _, err = srv.store.FindUserByEmail(ctx, req.GetEmail())

	// if err != nil {
	// 	if err != sql.ErrNoRows {
	// 		return nil, status.Errorf(codes.Internal, "error occured looking up the email")
	// 	}
	// }

	user, err := srv.store.CreateUser(ctx, db.CreateUserParams{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
		Password: hashedPw,
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

	// send verification email --

	return &pb.SignUpResponse{
		User: getUserResponse(user),
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
