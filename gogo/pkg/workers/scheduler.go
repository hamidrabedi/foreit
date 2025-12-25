package workers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

// Scheduler wraps Asynq scheduler for periodic tasks
type Scheduler struct {
	scheduler *asynq.Scheduler
	queue     *AsynqQueue
}

// NewScheduler creates a new scheduler using Asynq
func NewScheduler(queue *AsynqQueue) *Scheduler {
	scheduler := asynq.NewScheduler(queue.RedisConnOpt(), &asynq.SchedulerOpts{})

	return &Scheduler{
		scheduler: scheduler,
		queue:     queue,
	}
}

// Schedule schedules a recurring job using cron expression
func (s *Scheduler) Schedule(cronExpr string, job Job) error {
	task, err := s.queue.JobToTask(job)
	if err != nil {
		return err
	}

	entryID, err := s.scheduler.Register(cronExpr, task, asynq.Queue("default"))
	if err != nil {
		return err
	}

	log.Printf("Scheduled job %s with entry ID: %s", job.GetType(), entryID)
	return nil
}

// Start starts the scheduler
func (s *Scheduler) Start(ctx context.Context) error {
	return s.scheduler.Run()
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.scheduler.Shutdown()
}

// Schedule schedules a job with the default scheduler
func Schedule(cronExpr string, job Job) error {
	if defaultScheduler == nil {
		return fmt.Errorf("default scheduler not set")
	}
	return defaultScheduler.Schedule(cronExpr, job)
}

var defaultScheduler *Scheduler

// StartScheduler starts the default scheduler
func StartScheduler(ctx context.Context) error {
	if defaultQueue == nil {
		return fmt.Errorf("default queue not set")
	}

	asynqQueue, ok := defaultQueue.(*AsynqQueue)
	if !ok {
		return fmt.Errorf("default queue must be an AsynqQueue")
	}

	if defaultScheduler == nil {
		defaultScheduler = NewScheduler(asynqQueue)
	}

	go func() {
		if err := defaultScheduler.Start(ctx); err != nil {
			log.Printf("Scheduler error: %v", err)
		}
	}()

	return nil
}

// StopScheduler stops the default scheduler
func StopScheduler() {
	if defaultScheduler != nil {
		defaultScheduler.Stop()
	}
}
