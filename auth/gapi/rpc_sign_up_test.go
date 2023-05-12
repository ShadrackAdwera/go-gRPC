package gapi

import (
	"fmt"
	"testing"
	"time"

	mockdb "github.com/ShadrackAdwera/go-gRPC/db/mocks"
	db "github.com/ShadrackAdwera/go-gRPC/db/sqlc"
	"github.com/ShadrackAdwera/go-gRPC/pb"
	"github.com/ShadrackAdwera/go-gRPC/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

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
				store.EXPECT().CreateUser(gomock.Any(), gomock.Eq(db.CreateUserParams{
					Username: user.Username,
					Email:    user.Email,
					Password: password,
				})).Times(1).Return(gomock.Any(), nil)
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

		})
	}
}
