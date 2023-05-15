package db

import "context"

type CreateUserTxInput struct {
	CreateUserParams
	EmitCreateUser func(user User) error
}

type CreateUserTxOutput struct {
	User User `json:"user"`
}

func (s *Store) CreateUserTx(ctx context.Context, args CreateUserTxInput) (CreateUserTxOutput, error) {
	var output CreateUserTxOutput
	var err error

	err = s.execTx(ctx, func(q *Queries) error {
		output.User, err = q.CreateUser(ctx, CreateUserParams{
			Username: args.Username,
			Email:    args.Email,
			Password: args.Password,
		})

		if err != nil {
			return err
		}

		return args.EmitCreateUser(output.User)
	})

	return output, err
}
