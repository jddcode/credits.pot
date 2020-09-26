# Credits Pot

## Overview

The Credits Pot is a thread safe rate limiting mechanism which can be shared across multiple goroutines.

Imagine a pot of water. Each time you carry out a unit of work, you put a drop of water into your pot. When the pot
is full, no more work can be carried out. The pot has a leak and over time, water will leak out, allowing more work to
be completed.

Imagine your pot has a capacity of five and drips every two seconds. Your first five units of work will be
completed with no delay at all. After that your bucket will be full and work will pause. Every two seconds, a drop of
water will leak from your bucket and another unit of work can be completed.

Over the course of ten seconds, our example should process ten units of work; five initially and one every two seconds
for the remainder of the time.

## Example

```
func main() {

    // Create a new pot with a capacity of five, and which drips out a unit every two seconds
    pot := NewCreditsPot(CreditsPotConfig{ Size: 5, DripTime: time.Second * 2 })

    // Complete a unit of work whenever possible - should be five initially and then one every two seconds
    counter := 0
    for pot.Work() {

        fmt.Println(counter, "work units completed")
        counter++
    }    
}
```

## Core Library Alternative

There is a more comprehensive rate limiting package available in the `golang.org/x/time/rate` package, part of the
standard library. You can read the [documentation on that here](https://godoc.org/golang.org/x/time/rate).

The core library has considerably greater flexibility and functionality than this one, but is also more complex.