package credits

import (
	"errors"
	"sync"
	"time"
)

	func NewCreditsPot() CreditsPot {

		pot := creditsPot{}
		go pot.caretaker()
		return &pot
	}

	type CreditsPot interface {

		Work() error
	}

	type creditsPot struct {

		lock sync.RWMutex
		credits int
	}

	func (cp *creditsPot) Work() error {

		cp.lock.Lock()
		defer cp.lock.Unlock()

		if cp.credits > 0 {

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

			time.Sleep(time.Second)
		}
	}
