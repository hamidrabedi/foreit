package workers

import (
	"context"
	"fmt"
	"time"
)

type Job interface {
	Execute(ctx context.Context) error
	GetID() string
	GetType() string
	RetryCount() int
	MaxRetries() int
	IncrementRetry()
}

type BaseJob struct {
	ID             string
	Type           string
	CreatedAt      time.Time
	Retries        int
	MaxRetriesCount int
}

func NewBaseJob(jobType string) *BaseJob {
	return &BaseJob{
		ID:             generateID(),
		Type:           jobType,
		CreatedAt:      time.Now(),
		MaxRetriesCount: 3,
	}
}

func (j *BaseJob) GetID() string {
	return j.ID
}

func (j *BaseJob) GetType() string {
	return j.Type
}

func (j *BaseJob) RetryCount() int {
	return j.Retries
}

func (j *BaseJob) MaxRetries() int {
	return j.MaxRetriesCount
}

func (j *BaseJob) IncrementRetry() {
	j.Retries++
}

func (j *BaseJob) Execute(ctx context.Context) error {
	return fmt.Errorf("Execute method must be implemented by concrete job types")
}

type JobHandler interface {
	Handle(ctx context.Context, job Job) error
}

type JobRegistry struct {
	handlers map[string]JobHandler
}

var globalRegistry = &JobRegistry{
	handlers: make(map[string]JobHandler),
}

func Register(jobType string, handler JobHandler) {
	globalRegistry.handlers[jobType] = handler
}

func GetHandler(jobType string) (JobHandler, bool) {
	handler, ok := globalRegistry.handlers[jobType]
	return handler, ok
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
