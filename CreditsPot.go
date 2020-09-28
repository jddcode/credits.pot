package credits

import (
	"context"
	"errors"
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

		WaitFor(waitTime time.Duration) error
	}

	type creditsPot struct {

		lock sync.RWMutex
		credits []time.Time
		queueLock sync.RWMutex
		queue []struct{ RequestContext context.Context; ReadyToWork context.CancelFunc }
		config CreditsPotConfig
		nextExpiry time.Time
	}

	func (cp *creditsPot) WaitFor(waitTime time.Duration) error {

		requestCtx, _ := context.WithTimeout(context.Background(), waitTime)
		queueCtx := cp.joinQueue(requestCtx)

		for {

			select {

				case <- queueCtx.Done():

					return nil

				case <- requestCtx.Done():

					return errors.New("request_timeout")

				default:

					cp.iterate()
			}
			time.Sleep(time.Millisecond * 500)
		}
	}

	func (cp *creditsPot) joinQueue(requestCtx context.Context) context.Context {

		cp.queueLock.Lock()
		defer cp.queueLock.Unlock()

		// This context is Done when the caller is authorised to work.
		ctx, cancelFunc := context.WithCancel(context.Background())
		cp.queue = append(cp.queue, struct{ RequestContext context.Context; ReadyToWork context.CancelFunc }{ RequestContext: requestCtx, ReadyToWork: cancelFunc })
		return ctx
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
		for {

			cp.queueLock.Lock()
			if len(cp.queue) < 1 {

				cp.queueLock.Unlock()
				break
			}

			ticket := cp.queue[0]

			if len(cp.queue) == 1 {

				cp.queue = make([]struct{ RequestContext context.Context; ReadyToWork context.CancelFunc }, 0)

			} else {

				cp.queue = cp.queue[1:]
			}

			cp.queueLock.Unlock()

			// The signal that it is your turn is your ticket queue being closed
			workDone := false
			select {

				case <- ticket.RequestContext.Done():

					// This client has given up already. Move on
					continue

				default:

					// We can authorise this client to do their work
					ticket.ReadyToWork()
					workDone = true
			}

			if workDone {

				break
			}
		}

		return
	}

	type CreditsPotConfig struct {

		Size int
		DripTime time.Duration
	}
