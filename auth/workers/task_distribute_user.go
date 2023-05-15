package workers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type UserPayload struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

const TaskCreateUser = "task:create_user"

func (distro *TaskDistributor) DistributeUser(ctx context.Context, user UserPayload, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(user)

	if err != nil {
		return fmt.Errorf("error marshalling JSON %w", asynq.SkipRetry)
	}

	task := asynq.NewTask(TaskCreateUser, jsonPayload, opts...)

	info, err := distro.client.EnqueueContext(ctx, task)

	if err != nil {
		return fmt.Errorf("error occured %w", err)
	}

	log.Info().
		Str("task_type", info.Type).
		Str("task_id", info.ID).
		Str("queue", info.Queue).
		Bytes("payload", jsonPayload).
		Int("max_retries", info.MaxRetry).
		Msg("create user task enqueued")

	return nil
}
