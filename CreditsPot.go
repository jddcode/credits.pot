package credits

import (
	"sync"
	"time"
)

	func NewCreditsPot(config CreditsPotConfig) CreditsPot {

		pot := creditsPot{ credits: make([]time.Time, 0) }

		// Sort out the config
		if config.Size < 1 {

			config.Size = 1
		}

		if config.DripTime.Nanoseconds() < 1 {

			config.DripTime = time.Second
		}

		pot.config = config
		return &pot
	}

	type CreditsPot interface {

		Work()
	}

	type creditsPot struct {

		lock sync.RWMutex
		credits []time.Time
		queueLock sync.RWMutex
		queue []chan interface{}
		config CreditsPotConfig
		nextExpiry time.Time
	}

	func (cp *creditsPot) Work() {

		ticket := cp.joinQueue()
		for {

			select {

				case _, open := <- ticket:

					if !open {

						return
					}

				default:

					cp.iterate()
			}
			time.Sleep(time.Millisecond * 500)
		}
	}

	func (cp *creditsPot) joinQueue() chan interface{} {

		cp.queueLock.Lock()
		defer cp.queueLock.Unlock()

		ticket := make(chan interface{})
		cp.queue = append(cp.queue, ticket)
		return ticket
	}

	// This function does two things. Firstly it removes any credits which are expired, and then it checks to see if a new credit can be added.
	// If a new credit cannot be added (pot is full) it returns an error and the parent Work() function continues to wait.
	// If a new credit can be added, this is done - the new credit either has an expiry time of the configured expiry time after now, or that expiry time after the expiry time of the newest credit in the pot
	func (cp *creditsPot) iterate() {

		cp.lock.Lock()
		defer cp.lock.Unlock()

		newCredits := make([]time.Time, 0)
		for _, credit := range cp.credits {

			if time.Now().After(credit) {

				continue
			}

			newCredits = append(newCredits, credit)
		}

		cp.credits = newCredits

		if len(cp.credits) >= cp.config.Size {

			return
		}

		// If we don't have a next expiry yet or it's ages in the past because we haven't used the pot in a while, reset the nextExpiry
		if cp.nextExpiry.IsZero() || cp.nextExpiry.Add(cp.config.DripTime).Before(time.Now().Add(cp.config.DripTime)) {

			cp.nextExpiry = time.Now().Add(cp.config.DripTime)

		} else {

			cp.nextExpiry = cp.nextExpiry.Add(cp.config.DripTime)
		}

		cp.credits = append(cp.credits, cp.nextExpiry)

		// Notify the next in line they can do their work
		cp.queueLock.Lock()
		defer cp.queueLock.Unlock()

		if len(cp.queue) > 0 {

			ticket := cp.queue[0]

			if len(cp.queue) == 1 {

				cp.queue = make([]chan interface{}, 0)

			} else {

				cp.queue = cp.queue[1:]
			}

			// The signal that it is your turn is your ticket queue being closed
			close(ticket)
		}

		return
	}

	type CreditsPotConfig struct {

		Size int
		DripTime time.Duration
	}
