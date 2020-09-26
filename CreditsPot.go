package credits

import (
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

		Work()
	}

	type creditsPot struct {

		lock sync.RWMutex
		credits []time.Time
		config CreditsPotConfig
		nextExpiry time.Time
	}

	func (cp *creditsPot) Work() {

		for {

			err := cp.iterate()

			if err == nil {

				return
			}

			time.Sleep(time.Millisecond * 500)
		}
	}

	// This function does two things. Firstly it removes any credits which are expired, and then it checks to see if a new credit can be added.
	// If a new credit cannot be added (pot is full) it returns an error and the parent Work() function continues to wait.
	// If a new credit can be added, this is done - the new credit either has an expiry time of the configured expiry time after now, or that expiry time after the expiry time of the newest credit in the pot
	func (cp *creditsPot) iterate() error {

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

			return errors.New("over_limit")
		}

		// If we don't have a next expiry yet or it's ages in the past because we haven't used the pot in a while, reset the nextExpiry
		if cp.nextExpiry.IsZero() || cp.nextExpiry.Add(cp.config.DripTime).Before(time.Now().Add(cp.config.DripTime)) {

			cp.nextExpiry = time.Now().Add(cp.config.DripTime)

		} else {

			cp.nextExpiry = cp.nextExpiry.Add(cp.config.DripTime)
		}

		cp.credits = append(cp.credits, cp.nextExpiry)
		return nil
	}

	type CreditsPotConfig struct {

		Size int
		DripTime time.Duration
	}
