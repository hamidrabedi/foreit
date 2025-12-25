package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

type Worker struct {
	server      *asynq.Server
	mux         *asynq.ServeMux
	queue       *AsynqQueue
	concurrency int
}

func NewWorker(queue *AsynqQueue, concurrency int) *Worker {
	if concurrency <= 0 {
		concurrency = 10
	}

	redisOpt := queue.RedisConnOpt()
	config := asynq.Config{
		Concurrency: concurrency,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	}

	server := asynq.NewServer(redisOpt, config)
	mux := asynq.NewServeMux()

	for jobType, handler := range globalRegistry.handlers {
		jobType := jobType
		mux.HandleFunc(jobType, func(ctx context.Context, task *asynq.Task) error {
			job, err := taskToJob(task, jobType)
			if err != nil {
				return err
			}
			return handler.Handle(ctx, job)
		})
	}

	return &Worker{
		server:      server,
		mux:         mux,
		queue:       queue,
		concurrency: concurrency,
	}
}

func taskToJob(task *asynq.Task, jobType string) (Job, error) {
	var jobData map[string]interface{}
	if err := json.Unmarshal(task.Payload(), &jobData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	baseJob := &BaseJob{
		ID:             getString(jobData, "id", ""),
		Type:           jobType,
		MaxRetriesCount: 3,
	}

	if retries, ok := jobData["retries"].(float64); ok {
		baseJob.Retries = int(retries)
	}

	if maxRetries, ok := jobData["max_retries"].(float64); ok {
		baseJob.MaxRetriesCount = int(maxRetries)
	}

	if createdAtStr, ok := jobData["created_at"].(string); ok {
		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			baseJob.CreatedAt = t
		}
	}

	return baseJob, nil
}

func getString(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func (w *Worker) Start(ctx context.Context) error {
	log.Printf("Starting Asynq worker with concurrency: %d", w.concurrency)
	return w.server.Run(w.mux)
}

func (w *Worker) Stop() error {
	w.server.Shutdown()
	log.Println("Asynq worker stopped")
	return nil
}

func (w *Worker) IsRunning() bool {
	return w.server != nil
}

var defaultWorker *Worker

func Start(ctx context.Context, concurrency int) error {
	if defaultQueue == nil {
		return fmt.Errorf("default queue not set - call SetDefaultQueue first")
	}

	asynqQueue, ok := defaultQueue.(*AsynqQueue)
	if !ok {
		return fmt.Errorf("default queue must be an AsynqQueue")
	}

	if defaultWorker == nil {
		defaultWorker = NewWorker(asynqQueue, concurrency)
	}

	go func() {
		if err := defaultWorker.Start(ctx); err != nil {
			log.Printf("Worker error: %v", err)
		}
	}()

	return nil
}

func Stop() error {
	if defaultWorker != nil {
		return defaultWorker.Stop()
	}
	return nil
}
