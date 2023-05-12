package gapi

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	mockdb "github.com/ShadrackAdwera/go-gRPC/db/mocks"
	db "github.com/ShadrackAdwera/go-gRPC/db/sqlc"
	"github.com/ShadrackAdwera/go-gRPC/pb"
	"github.com/ShadrackAdwera/go-gRPC/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
	user     db.User
}

func (expected eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	actualArg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := utils.IsPassword(expected.password, actualArg.Password)
	if err != nil {
		return false
	}

	expected.arg.Password = actualArg.Password
	return reflect.DeepEqual(expected.arg, actualArg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string, user db.User) gomock.Matcher {
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
		buildStubs func(store *mockdb.MockTxStore)
		comparator func(t *testing.T, res *pb.SignUpResponse, err error)
	}{
		{
			name: "TestOK",
			body: &pb.SignUpRequest{
				Username: user.Username,
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockTxStore) {
				args := db.CreateUserParams{
					Username: user.Username,
					Email:    user.Email,
					Password: password,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(args, password, user)).Times(1).Return(user, nil)
			},
			comparator: func(t *testing.T, res *pb.SignUpResponse, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res.User)
				pbUser := res.GetUser()

				require.Equal(t, user.Username, pbUser.Username)
				require.Equal(t, user.Email, pbUser.Email)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctlr := gomock.NewController(t)
			store := mockdb.NewMockTxStore(ctlr)

			defer ctlr.Finish()

			testCase.buildStubs(store)

			srv := newTestServer(t, store)
			res, err := srv.SignUp(context.Background(), testCase.body)
			testCase.comparator(t, res, err)
		})
	}
}
