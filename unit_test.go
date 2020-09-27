package credits

import (
	"testing"
	"time"
)

var pot CreditsPot
var workerA int
var workerB int

	func TestRestrictions(t *testing.T) {

		pot = NewCreditsPot(CreditsPotConfig{ Size: 5, DripTime: time.Second * 2 })

		end := time.Now().Add(time.Second * 10)
		go worker(&workerA, end)
		go worker(&workerB, end)

		for time.Now().Before(end.Add(time.Second)) {

			time.Sleep(time.Second)
		}

		completedWork := workerA + workerB
		if completedWork != 10 {

			t.Error("Completed", completedWork, "units of work. Should only be able to complete ten")
		}
	}

	func worker(workStack *int, end time.Time) {

		for time.Now().Before(end) {

			pot.Work()
			*workStack++
		}
	}
