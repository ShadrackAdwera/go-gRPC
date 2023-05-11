package token

import (
	"fmt"
	"time"

	uuid "github.com/google/uuid"
)

var (
	ErrorExpiredToken = fmt.Errorf("your token is expired")
)

type TokenPayload struct {
	TokenId   uuid.UUID `json:"token_id"`
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(id int64, username string, email string, duration time.Duration) (*TokenPayload, error) {
	tokenId, err := uuid.NewRandom()

	if err != nil {
		return &TokenPayload{}, nil
	}

	payload := &TokenPayload{
		TokenId:   tokenId,
		ID:        id,
		Username:  username,
		Email:     email,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *TokenPayload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrorExpiredToken
	}
	return nil
}
