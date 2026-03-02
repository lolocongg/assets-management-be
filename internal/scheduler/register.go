package scheduler

import "github.com/davidcm146/assets-management-be.git/internal/scheduler/jobs"

func RegisterJobs(s *Scheduler, jobs []jobs.Job) error {
	for _, job := range jobs {
		err := s.Add(job.Schedule(), job.Run)
		if err != nil {
			return err
		}
	}
	return nil
}
