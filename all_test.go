package credits

import (
	"testing"
	"time"
)

	func TestRestrictions(t *testing.T) {

		myPot := NewCreditsPot()

		start := time.Now()
		score := 0
		for time.Now().Before(start.Add(time.Second * 2)) {

			err := myPot.Work()

			if err == nil {

				score++
			}
		}

		if score != 3 {

			t.Error("Completed", score, "units of work. Should only be able to complete three")
		}
	}
