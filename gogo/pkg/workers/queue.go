package workers

import (
	"context"
	"fmt"
	"time"
)

// Queue represents a job queue
type Queue interface {
	// Enqueue adds a job to the queue
	Enqueue(ctx context.Context, job Job) error
	
	// EnqueueDelayed adds a job with a delay
	EnqueueDelayed(ctx context.Context, job Job, delay time.Duration) error
}

// AsynqConfig configures Asynq queue
type AsynqConfig struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	Concurrency   int
}

// DefaultAsynqConfig returns default Asynq configuration
func DefaultAsynqConfig() *AsynqConfig {
	return &AsynqConfig{
		RedisAddr:     "localhost:6379",
		RedisPassword: "",
		RedisDB:       0,
		Concurrency:   10,
	}
}

// NewAsynqQueueFromConfig creates an Asynq queue from config
func NewAsynqQueueFromConfig(config *AsynqConfig) (Queue, error) {
	return NewAsynqQueue(config.RedisAddr, config.RedisPassword, config.RedisDB)
}

// Enqueue adds a job to the default queue
func Enqueue(ctx context.Context, job Job) error {
	if defaultQueue == nil {
		return fmt.Errorf("default queue not set")
	}
	return defaultQueue.Enqueue(ctx, job)
}

// EnqueueDelayed adds a job to the queue with a delay
func EnqueueDelayed(ctx context.Context, job Job, delay time.Duration) error {
	if defaultQueue == nil {
		return fmt.Errorf("default queue not set")
	}
	return defaultQueue.EnqueueDelayed(ctx, job, delay)
}

var defaultQueue Queue

// SetDefaultQueue sets the default queue
func SetDefaultQueue(queue Queue) {
	defaultQueue = queue
}

// GetDefaultQueue returns the default queue
func GetDefaultQueue() Queue {
	return defaultQueue
}
