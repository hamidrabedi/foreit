package workers

import (
	"context"
	"fmt"
	"time"
)

type Queue interface {
	Enqueue(ctx context.Context, job Job) error
	EnqueueDelayed(ctx context.Context, job Job, delay time.Duration) error
}

type AsynqConfig struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	Concurrency   int
}

func DefaultAsynqConfig() *AsynqConfig {
	return &AsynqConfig{
		RedisAddr:     "localhost:6379",
		RedisPassword: "",
		RedisDB:       0,
		Concurrency:   10,
	}
}

func NewAsynqQueueFromConfig(config *AsynqConfig) (Queue, error) {
	return NewAsynqQueue(config.RedisAddr, config.RedisPassword, config.RedisDB)
}

func Enqueue(ctx context.Context, job Job) error {
	if defaultQueue == nil {
		return fmt.Errorf("default queue not set")
	}
	return defaultQueue.Enqueue(ctx, job)
}

func EnqueueDelayed(ctx context.Context, job Job, delay time.Duration) error {
	if defaultQueue == nil {
		return fmt.Errorf("default queue not set")
	}
	return defaultQueue.EnqueueDelayed(ctx, job, delay)
}

var defaultQueue Queue

func SetDefaultQueue(queue Queue) {
	defaultQueue = queue
}

func GetDefaultQueue() Queue {
	return defaultQueue
}
