package mapbased

import "time"

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor) Run(s *Storage) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			s.DeleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func stopJanitor(s *Storage) {
	s.janitor.stop <- true
}

func runJanitor(s *Storage, ci time.Duration) {
	j := &janitor{
		Interval: ci,
		stop:     make(chan bool),
	}
	s.janitor = j
	go j.Run(s)
}
