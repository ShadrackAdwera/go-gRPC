package token

import "time"

type TokenMaker interface {
	CreateToken(username string, id int64, email string, duration time.Duration) (*TokenPayload, string, error)
	VerifyToken(token string) (*TokenPayload, error)
}
