package scheduler

import (
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
}

func NewScheduler() *Scheduler {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil || loc == nil {
		loc = time.FixedZone("ICT", 7*3600)
	}

	return &Scheduler{
		cron: cron.New(
			cron.WithLocation(loc),
			cron.WithChain(
				cron.Recover(cron.DefaultLogger),
				cron.SkipIfStillRunning(cron.DefaultLogger),
			),
		),
	}
}

func (s *Scheduler) Add(spec string, job func()) error {
	_, err := s.cron.AddFunc(spec, job)
	return err
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}
