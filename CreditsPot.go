package credits

import (
	"errors"
	"sync"
	"time"
)

	func NewCreditsPot(config CreditsPotConfig) CreditsPot {

		pot := creditsPot{}

		// Sort out the config
		if config.Burst < 1 {

			config.Burst = 1
		}

		if config.DecrementSeconds < 1 {

			config.DecrementSeconds = 1
		}

		pot.config = config
		go pot.caretaker()
		return &pot
	}

	type CreditsPot interface {

		Work() error
	}

	type creditsPot struct {

		lock sync.RWMutex
		credits int
		config CreditsPotConfig
	}

	func (cp *creditsPot) Work() error {

		cp.lock.Lock()
		defer cp.lock.Unlock()

		if cp.credits >= cp.config.Burst {

			return errors.New("over_limit")
		}

		cp.credits++
		return nil
	}

	func (cp *creditsPot) caretaker() {

		for {

			cp.lock.RLock()
			amount := cp.credits
			cp.lock.RUnlock()

			if amount > 0 {

				cp.lock.Lock()
				cp.credits--
				cp.lock.Unlock()
			}

			time.Sleep(time.Second * time.Duration(cp.config.DecrementSeconds))
		}
	}

	type CreditsPotConfig struct {

		Burst int
		DecrementSeconds int
	}
