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

			err := pot.WaitFor(time.Minute)

			if err == nil {

				*workStack++
			}
		}
	}

	func TestWorkTimeout(t *testing.T) {

		// Create a very slow pot to test that requests time out correctly
		pot = NewCreditsPot(CreditsPotConfig{ Size: 1, DripTime: time.Minute * 10 })

		// Do one piece of work
		err := pot.WaitFor(time.Minute)

		if err != nil {

			t.Error("First item into pot had an error but should have succeeded :", err.Error())
			return
		}

		// Now this should time out
		err = pot.WaitFor(time.Second)

		if err == nil {

			t.Error("Second item into pot had no error but it should have timed out")
		}

		if err.Error() != "request_timeout" {

			t.Error("Second item into pot should have timed out it's request, but it gave a different error code :", err.Error())
		}
	}
