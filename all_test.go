package credits

import (
	"fmt"
	"testing"
	"time"
)

	func TestRestrictions(t *testing.T) {

		myPot := NewCreditsPot(CreditsPotConfig{ Size: 5, DripSeconds: 2 })

		start := time.Now()
		score := 0
		for time.Now().Before(start.Add(time.Second * 10)) {

			if !myPot.Work() {

				t.Error("Test aborted")
				return
			}

			score++
		}

		fmt.Println(start, time.Now())

		if score != 10 {

			t.Error("Completed", score, "units of work. Should only be able to complete ten")
		}
	}
