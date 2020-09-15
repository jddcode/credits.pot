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

		if config.DripSeconds < 1 {

			config.DripSeconds = 1
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
		}
	}

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
		if cp.nextExpiry.IsZero() || cp.nextExpiry.Add(time.Second).Before(time.Now().Add(time.Second * time.Duration(cp.config.DripSeconds))) {

			cp.nextExpiry = time.Now().Add(time.Second * time.Duration(cp.config.DripSeconds))

		} else {

			cp.nextExpiry = cp.nextExpiry.Add(time.Second * time.Duration(cp.config.DripSeconds))
		}

		cp.credits = append(cp.credits, cp.nextExpiry)
		return nil
	}

	type CreditsPotConfig struct {

		Size int
		DripSeconds int
	}
