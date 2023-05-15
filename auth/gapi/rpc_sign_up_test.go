package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	mockdb "github.com/ShadrackAdwera/go-gRPC/db/mocks"
	db "github.com/ShadrackAdwera/go-gRPC/db/sqlc"
	"github.com/ShadrackAdwera/go-gRPC/pb"
	"github.com/ShadrackAdwera/go-gRPC/utils"
	"github.com/ShadrackAdwera/go-gRPC/workers"
	mockworkers "github.com/ShadrackAdwera/go-gRPC/workers/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserTxInput
	password string
	user     db.User
}

func (expected eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	actualArg, ok := x.(db.CreateUserTxInput)
	if !ok {
		return false
	}

	err := utils.IsPassword(expected.password, actualArg.Password)
	if err != nil {
		return false
	}

	expected.arg.Password = actualArg.Password

	if !reflect.DeepEqual(expected.arg.CreateUserParams, actualArg.CreateUserParams) {
		return false
	}

	err = actualArg.EmitCreateUser(expected.user)

	return err == nil
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserTxInput, password string, user db.User) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password, user}
}

func createTestUser() (db.User, string) {
	username := utils.RandomString(12)
	password := utils.RandomString(16)
	hashPw, _ := utils.HashPassword(password)

	return db.User{
		ID:        utils.RandomInteger(1, 100),
		Username:  username,
		Email:     fmt.Sprintf("%s@mail.com", username),
		Password:  hashPw,
		CreatedAt: time.Now(),
	}, password
}

func TestSignUp(t *testing.T) {
	user, password := createTestUser()

	testCases := []struct {
		name       string
		body       *pb.SignUpRequest
		buildStubs func(store *mockdb.MockTxStore, distro *mockworkers.MockDistributor)
		comparator func(t *testing.T, res *pb.SignUpResponse, err error)
	}{
		{
			name: "TestOK",
			body: &pb.SignUpRequest{
				Username: user.Username,
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockTxStore, distro *mockworkers.MockDistributor) {
				args := db.CreateUserTxInput{
					CreateUserParams: db.CreateUserParams{
						Username: user.Username,
						Email:    user.Email,
						Password: password,
					},
				}

				distroPayload := workers.UserPayload{
					ID:       user.ID,
					Username: user.Username,
					Email:    user.Email,
				}
				store.EXPECT().CreateUserTx(gomock.Any(), EqCreateUserParams(args, password, user)).Times(1).Return(db.CreateUserTxOutput{
					User: user,
				}, nil)
				distro.EXPECT().DistributeUser(gomock.Any(), distroPayload, gomock.Any()).Times(1).Return(nil)
			},
			comparator: func(t *testing.T, res *pb.SignUpResponse, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res.User)
				pbUser := res.GetUser()

				require.Equal(t, user.Username, pbUser.Username)
				require.Equal(t, user.Email, pbUser.Email)
			},
		},
		{
			name: "TestInternalServerError",
			body: &pb.SignUpRequest{
				Username: user.Username,
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockTxStore, distro *mockworkers.MockDistributor) {
				store.EXPECT().CreateUserTx(gomock.Any(), gomock.Any()).Times(1).Return(db.CreateUserTxOutput{}, sql.ErrConnDone)
				distro.EXPECT().DistributeUser(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			comparator: func(t *testing.T, res *pb.SignUpResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			storeCtlr := gomock.NewController(t)
			defer storeCtlr.Finish()
			store := mockdb.NewMockTxStore(storeCtlr)

			workerCtlr := gomock.NewController(t)
			defer workerCtlr.Finish()
			distro := mockworkers.NewMockDistributor(workerCtlr)

			testCase.buildStubs(store, distro)

			srv := newTestServer(t, store, distro)
			res, err := srv.SignUp(context.Background(), testCase.body)
			testCase.comparator(t, res, err)
		})
	}
}
