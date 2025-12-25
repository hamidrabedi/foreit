package workers

import (
	"context"
	"encoding/json"
	"time"
)

// Job represents a background job
type Job interface {
	// Execute executes the job
	Execute(ctx context.Context) error
	
	// GetID returns the job ID
	GetID() string
	
	// GetType returns the job type
	GetType() string
	
	// RetryCount returns the current retry count
	RetryCount() int
	
	// MaxRetries returns the maximum number of retries
	MaxRetries() int
	
	// IncrementRetry increments the retry count
	IncrementRetry()
}

// BaseJob provides base job functionality
type BaseJob struct {
	ID        string
	Type      string
	CreatedAt time.Time
	Retries   int
	MaxRetries int
}

// NewBaseJob creates a new base job
func NewBaseJob(jobType string) *BaseJob {
	return &BaseJob{
		ID:        generateID(),
		Type:      jobType,
		CreatedAt: time.Now(),
		MaxRetries: 3,
	}
}

// GetID returns the job ID
func (j *BaseJob) GetID() string {
	return j.ID
}

// GetType returns the job type
func (j *BaseJob) GetType() string {
	return j.Type
}

// RetryCount returns the current retry count
func (j *BaseJob) RetryCount() int {
	return j.Retries
}

// MaxRetries returns the maximum number of retries
func (j *BaseJob) MaxRetries() int {
	return j.MaxRetries
}

// IncrementRetry increments the retry count
func (j *BaseJob) IncrementRetry() {
	j.Retries++
}

// JobHandler handles a specific job type
type JobHandler interface {
	Handle(ctx context.Context, job Job) error
}

// JobRegistry stores job handlers
type JobRegistry struct {
	handlers map[string]JobHandler
}

var globalRegistry = &JobRegistry{
	handlers: make(map[string]JobHandler),
}

// Register registers a job handler
func Register(jobType string, handler JobHandler) {
	globalRegistry.handlers[jobType] = handler
}

// GetHandler retrieves a job handler
func GetHandler(jobType string) (JobHandler, bool) {
	handler, ok := globalRegistry.handlers[jobType]
	return handler, ok
}

// generateID generates a unique job ID
func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string (simplified)
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// JobData represents job data for serialization
type JobData struct {
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
	ID        string          `json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	Retries   int             `json:"retries"`
	MaxRetries int            `json:"max_retries"`
}

// Serialize serializes a job to JSON
func Serialize(job Job) ([]byte, error) {
	data, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}
	
	jobData := JobData{
		Type:      job.GetType(),
		Data:      data,
		ID:        job.GetID(),
		CreatedAt: time.Now(),
	}
	
	return json.Marshal(jobData)
}

// Deserialize deserializes a job from JSON
func Deserialize(data []byte, factory func(string) Job) (Job, error) {
	var jobData JobData
	if err := json.Unmarshal(data, &jobData); err != nil {
		return nil, err
	}
	
	return factory(jobData.Type), nil
}

