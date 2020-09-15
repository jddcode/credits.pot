package credits

import (
	"testing"
	"time"
)

	func TestRestrictions(t *testing.T) {

		myPot := NewCreditsPot(CreditsPotConfig{ Burst: 5, DecrementSeconds: 2 })

		start := time.Now()
		score := 0
		for time.Now().Before(start.Add(time.Second * 10)) {

			err := myPot.Work()

			if err == nil {

				score++
			}
		}

		if score != 10 {

			t.Error("Completed", score, "units of work. Should only be able to complete ten")
		}
	}
