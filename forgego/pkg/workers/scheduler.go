package workers

import (
	"context"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

type Scheduler struct {
	scheduler *asynq.Scheduler
	queue     *AsynqQueue
}

func NewScheduler(queue *AsynqQueue) *Scheduler {
	scheduler := asynq.NewScheduler(queue.RedisConnOpt(), &asynq.SchedulerOpts{})

	return &Scheduler{
		scheduler: scheduler,
		queue:     queue,
	}
}

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

func (s *Scheduler) Start(ctx context.Context) error {
	return s.scheduler.Run()
}

func (s *Scheduler) Stop() {
	s.scheduler.Shutdown()
}

func Schedule(cronExpr string, job Job) error {
	if defaultScheduler == nil {
		return fmt.Errorf("default scheduler not set")
	}
	return defaultScheduler.Schedule(cronExpr, job)
}

var defaultScheduler *Scheduler

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

func StopScheduler() {
	if defaultScheduler != nil {
		defaultScheduler.Stop()
	}
}
