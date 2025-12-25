# Workers Module

Background job processing powered by [Asynq](https://github.com/hibiken/asynq) - a simple, reliable, and efficient distributed task queue for Go.

## Features

- ✅ **Redis-based** - Uses Redis (already used for caching)
- ✅ **Task Retries** - Automatic retries with exponential backoff
- ✅ **Scheduled Tasks** - Cron-like scheduling
- ✅ **Priority Queues** - Multiple queues with different priorities
- ✅ **Web UI** - Built-in monitoring dashboard
- ✅ **Rate Limiting** - Built-in rate limiting support
- ✅ **Active Development** - Regularly updated and maintained

## Installation

Asynq is automatically included when you use the workers module.

## Usage

### Setup

```go
import (
    "github.com/forgego/forge/pkg/workers"
)

// Create Asynq queue
queue, err := workers.NewAsynqQueue("localhost:6379", "", 0)
if err != nil {
    log.Fatal(err)
}

// Set as default queue
workers.SetDefaultQueue(queue)
```

### Define a Job

```go
type SendEmailJob struct {
    workers.BaseJob
    To      string
    Subject string
    Body    string
}

func (j *SendEmailJob) Execute(ctx context.Context) error {
    return sendEmail(j.To, j.Subject, j.Body)
}
```

### Register Job Handler

```go
workers.Register("send_email", &EmailJobHandler{})

type EmailJobHandler struct{}

func (h *EmailJobHandler) Handle(ctx context.Context, job workers.Job) error {
    emailJob := job.(*SendEmailJob)
    return emailJob.Execute(ctx)
}
```

### Enqueue Jobs

```go
job := &SendEmailJob{
    BaseJob: *workers.NewBaseJob("send_email"),
    To: "user@example.com",
    Subject: "Welcome",
    Body: "Welcome to our app!",
}

// Enqueue immediately
workers.Enqueue(ctx, job)

// Enqueue with delay
workers.EnqueueDelayed(ctx, job, 5*time.Minute)
```

### Start Workers

```go
// Start with 10 workers
workers.Start(ctx, 10)

// Stop workers
defer workers.Stop()
```

### Scheduled Jobs (Cron)

```go
// Schedule a daily job (runs at 6 AM every day)
workers.Schedule("0 6 * * *", &DailyReportJob{})

// Schedule every hour
workers.Schedule("0 * * * *", &HourlyCleanupJob{})

// Start scheduler
workers.StartScheduler(ctx)
```

### Web UI

Asynq provides a built-in web UI for monitoring tasks. To enable it:

```go
import (
    "github.com/hibiken/asynqmon"
    "github.com/forgego/forge/pkg/workers"
)

// Get the Asynq queue
queue := workers.GetDefaultQueue().(*workers.AsynqQueue)

// Create monitor
monitor := asynqmon.New(asynqmon.Options{
    RootPath:     "/monitor",
    RedisConnOpt: queue.RedisConnOpt(),
})

// Mount to your Fiber app
app.Use("/monitor", monitor)
```

## Queue Priorities

Asynq supports multiple queues with different priorities:

- `critical` - High priority tasks (weight: 6)
- `default` - Normal priority tasks (weight: 3)
- `low` - Low priority tasks (weight: 1)

## Benefits of Asynq

1. **Simple API** - Easy to use and integrate
2. **Redis-based** - Leverages existing Redis infrastructure
3. **Active Development** - Regularly updated with new features
4. **Built-in Monitoring** - Web UI for task management
5. **Production Ready** - Used by many companies in production
6. **Better Performance** - Optimized for high-throughput scenarios

## Migration from Custom Queue

If you were using the old in-memory queue, simply replace:

```go
// Old
queue := workers.NewMemoryQueue()
workers.SetDefaultQueue(queue)

// New
queue, _ := workers.NewAsynqQueue("localhost:6379", "", 0)
workers.SetDefaultQueue(queue)
```

The rest of your code remains the same!
