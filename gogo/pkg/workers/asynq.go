package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// AsynqQueue wraps Asynq client and server
type AsynqQueue struct {
	client   *asynq.Client
	server   *asynq.Server
	mux      *asynq.ServeMux
	redisOpt asynq.RedisClientOpt
}

// NewAsynqQueue creates a new Asynq-based queue
func NewAsynqQueue(redisAddr string, redisPassword string, redisDB int) (*AsynqQueue, error) {
	redisOpt := asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	}

	client := asynq.NewClient(redisOpt)
	server := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	})

	mux := asynq.NewServeMux()

	return &AsynqQueue{
		client:   client,
		server:   server,
		mux:      mux,
		redisOpt: redisOpt,
	}, nil
}

// Enqueue adds a job to the queue
func (q *AsynqQueue) Enqueue(ctx context.Context, job Job) error {
	task, err := q.jobToTask(job)
	if err != nil {
		return err
	}

	opts := []asynq.Option{
		asynq.MaxRetry(job.MaxRetries()),
		asynq.Timeout(30 * time.Minute),
	}

	if job.RetryCount() > 0 {
		opts = append(opts, asynq.ProcessIn(time.Duration(job.RetryCount())*time.Second))
	}

	_, err = q.client.Enqueue(task, opts...)
	return err
}

// EnqueueDelayed adds a job with a delay
func (q *AsynqQueue) EnqueueDelayed(ctx context.Context, job Job, delay time.Duration) error {
	task, err := q.jobToTask(job)
	if err != nil {
		return err
	}

	_, err = q.client.Enqueue(task, asynq.ProcessIn(delay))
	return err
}

// jobToTask converts our Job interface to Asynq task
func (q *AsynqQueue) jobToTask(job Job) (*asynq.Task, error) {
	data, err := json.Marshal(job)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal job: %w", err)
	}

	return asynq.NewTask(job.GetType(), data), nil
}

// RegisterHandler registers a job handler
func (q *AsynqQueue) RegisterHandler(jobType string, handler func(context.Context, *asynq.Task) error) {
	q.mux.HandleFunc(jobType, handler)
}

// Start starts the Asynq server
func (q *AsynqQueue) Start() error {
	return q.server.Run(q.mux)
}

// Shutdown gracefully shuts down the server
func (q *AsynqQueue) Shutdown() {
	q.server.Shutdown()
	q.client.Close()
}

// GetClient returns the Asynq client
func (q *AsynqQueue) GetClient() *asynq.Client {
	return q.client
}

// GetServer returns the Asynq server
func (q *AsynqQueue) GetServer() *asynq.Server {
	return q.server
}

// RedisConnOpt returns the Redis connection options
func (q *AsynqQueue) RedisConnOpt() asynq.RedisClientOpt {
	return q.redisOpt
}

// JobToTask converts our Job interface to Asynq task (exported)
func (q *AsynqQueue) JobToTask(job Job) (*asynq.Task, error) {
	return q.jobToTask(job)
}

