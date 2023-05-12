package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/ShadrackAdwera/go-gRPC/db/sqlc"
	"github.com/ShadrackAdwera/go-gRPC/pb"
	"github.com/ShadrackAdwera/go-gRPC/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (srv *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// check if account exists
	user, err := srv.store.FindUserByEmail(ctx, req.GetEmail())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "account does not exist %s", err)
		}
		return nil, status.Errorf(codes.Internal, "an error occured while logging in %s", err)
	}

	// check password
	err = utils.IsPassword(req.GetPassword(), user.Password)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "auth failed %s", err)
	}

	// create access token
	aPayload, aTkn, err := srv.maker.CreateToken(user.Username, user.ID, user.Email, time.Minute*15)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "an error occured %s", err)
	}

	// create refresh token
	rPayload, rTkn, err := srv.maker.CreateToken(user.Username, user.ID, user.Email, time.Hour)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "an error occured %s", err)
	}

	// extract metadata from incoming request
	meta := extractMetadata(ctx)

	// create session
	sess, err := srv.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           rPayload.TokenId,
		Username:     rPayload.Username,
		UserID:       rPayload.ID,
		RefreshToken: rTkn,
		UserAgent:    meta.UserAgent,
		ClientIp:     meta.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    rPayload.ExpiredAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "an error occured creating the session %s", err)
	}

	response := &pb.LoginResponse{
		User:                       getUserResponse(user),
		SessionId:                  sess.ID.String(),
		AccessToken:                aTkn,
		RefreshToken:               sess.RefreshToken,
		AccessTokenExpirationTime:  timestamppb.New(aPayload.ExpiredAt),
		RefreshTokenExpirationTime: timestamppb.New(rPayload.ExpiredAt),
	}
	return response, nil
}
