package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (TokenMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("provided symmetric key of length: %d is not equal to required key length: %d", len(symmetricKey), chacha20poly1305.KeySize)
	}
	p := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return p, nil
}

func (p *PasetoMaker) CreateToken(username string, id int64, email string, duration time.Duration) (*TokenPayload, string, error) {
	payload, err := NewPayload(id, username, email, duration)
	if err != nil {
		return nil, "", err
	}

	tkn, err := p.paseto.Encrypt(p.symmetricKey, payload, nil)

	if err != nil {
		return nil, "", err
	}

	return payload, tkn, nil
}

func (p *PasetoMaker) VerifyToken(token string) (*TokenPayload, error) {

	tknPayload := &TokenPayload{}

	err := p.paseto.Decrypt(token, p.symmetricKey, tknPayload, nil)

	if err != nil {
		return nil, err
	}

	err = tknPayload.Valid()

	if err != nil {
		return nil, err
	}

	return tknPayload, nil
}
