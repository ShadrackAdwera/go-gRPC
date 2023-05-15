package workers

import (
	"context"

	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type Distributor interface {
	DistributeUser(ctx context.Context, user UserPayload, opts ...asynq.Option) error
}

type TaskDistributor struct {
	client *asynq.Client
}

func NewDistributor(opts *asynq.RedisClientOpt) Distributor {
	client := asynq.NewClient(opts)

	return &TaskDistributor{
		client,
	}
}
