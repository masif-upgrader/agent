package main

import (
	"sync"
	"time"
)

type doneActions struct {
	amount   uint64
	from, to time.Time
}

type statsBookkeeper struct {
	sync.RWMutex

	history []doneActions
}

func (s *statsBookkeeper) addDoneActions(amount uint64, from, to time.Time) {
	s.Lock()
	defer s.Unlock()

	threshold := time.Now().Add(-time.Minute * 15)
	i := 0

	for i < len(s.history) && !s.history[i].to.After(threshold) {
		i++
	}

	if i > 0 {
		s.history = s.history[i:]
	}

	s.history = append(s.history, doneActions{amount, from, to})
}

func (s *statsBookkeeper) queryLoad() [3]float64 {
	type doneActionsPerPeriod struct {
		amount   float64
		from, to time.Time
	}

	const res = time.Microsecond

	s.RLock()
	history := s.history
	s.RUnlock()

	now := time.Now()

	load := [3]doneActionsPerPeriod{
		{0.0, now.Add(-time.Minute), now},
		{0.0, now.Add(-time.Minute * 5), now},
		{0.0, now.Add(-time.Minute * 15), now},
	}

	for _, entry := range history {
		for i := range load {
			if !(entry.to.Before(load[i].from) || entry.from.After(load[i].to)) {
				div := float64(entry.to.Sub(entry.from) / res)

				if div == 0.0 {
					load[i].amount += float64(entry.amount)
				} else {
					from := entry.from
					if from.Before(load[i].from) {
						from = load[i].from
					}

					to := entry.to
					if to.After(load[i].to) {
						to = load[i].to
					}

					load[i].amount += float64(entry.amount) * (float64(to.Sub(from)/res) / div)
				}
			}
		}
	}

	return [3]float64{load[0].amount, load[1].amount, load[2].amount}
}
