package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

type AsynqQueue struct {
	client   *asynq.Client
	server   *asynq.Server
	mux      *asynq.ServeMux
	redisOpt asynq.RedisClientOpt
}

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

func (q *AsynqQueue) EnqueueDelayed(ctx context.Context, job Job, delay time.Duration) error {
	task, err := q.jobToTask(job)
	if err != nil {
		return err
	}

	_, err = q.client.Enqueue(task, asynq.ProcessIn(delay))
	return err
}

func (q *AsynqQueue) jobToTask(job Job) (*asynq.Task, error) {
	data, err := json.Marshal(job)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal job: %w", err)
	}

	return asynq.NewTask(job.GetType(), data), nil
}

func (q *AsynqQueue) RegisterHandler(jobType string, handler func(context.Context, *asynq.Task) error) {
	q.mux.HandleFunc(jobType, handler)
}

func (q *AsynqQueue) Start() error {
	return q.server.Run(q.mux)
}

func (q *AsynqQueue) Shutdown() {
	q.server.Shutdown()
	q.client.Close()
}

func (q *AsynqQueue) GetClient() *asynq.Client {
	return q.client
}

func (q *AsynqQueue) GetServer() *asynq.Server {
	return q.server
}

func (q *AsynqQueue) RedisConnOpt() asynq.RedisClientOpt {
	return q.redisOpt
}

func (q *AsynqQueue) JobToTask(job Job) (*asynq.Task, error) {
	return q.jobToTask(job)
}
