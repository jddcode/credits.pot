package credits

import (
	"errors"
	"sync"
	"time"
)

	func NewCreditsPot(config CreditsPotConfig) CreditsPot {

		pot := creditsPot{ credits: make([]time.Time, 0) }

		// Sort out the config
		if config.Burst < 1 {

			config.Burst = 1
		}

		if config.DecrementSeconds < 1 {

			config.DecrementSeconds = 1
		}

		pot.config = config
		return &pot
	}

	type CreditsPot interface {

		Work() error
	}

	type creditsPot struct {

		lock sync.RWMutex
		credits []time.Time
		config CreditsPotConfig
		nextExpiry time.Time
	}

	func (cp *creditsPot) Work() error {

		cp.lock.Lock()
		defer cp.lock.Unlock()

		newCredits := make([]time.Time, 0)
		for _, credit := range cp.credits {

			if time.Now().Format("05") == credit.Format("05") || time.Now().After(credit) {

				continue
			}

			newCredits = append(newCredits, credit)
		}

		cp.credits = newCredits

		if len(cp.credits) >= cp.config.Burst {

			return errors.New("over_limit")
		}

		if cp.nextExpiry.IsZero() {

			cp.nextExpiry = time.Now().Add(time.Second * time.Duration(cp.config.DecrementSeconds))

		} else {

			cp.nextExpiry = cp.nextExpiry.Add(time.Second * time.Duration(cp.config.DecrementSeconds))
		}

		cp.credits = append(cp.credits, cp.nextExpiry)
		return nil
	}

	type CreditsPotConfig struct {

		Burst int
		DecrementSeconds int
	}
