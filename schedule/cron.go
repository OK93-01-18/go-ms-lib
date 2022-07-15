package schedule

import robficCron "github.com/robfig/cron/v3"

type Croner interface {
	Start()
	AddFunc(string, func()) error
}

type Cron struct {
	cron *robficCron.Cron
}

func (s *Cron) Start() {
	s.cron.Start()
}

func (s *Cron) AddFunc(spec string, f func()) error {
	_, err := s.cron.AddFunc(spec, f)
	return err
}

func NewCron() Croner {
	return &Cron{cron: robficCron.New(robficCron.WithSeconds())}
}
