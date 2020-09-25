# Credits Pot

## Overview

The Credits Pot is a rate limiting mechanism. To use, you set the size of your pot and the number of second between
each drip of water from your pot.

Example: Imagine you have a pot with a capacity of `ten` and a drip seconds value of `2`. You will initially be able
to carry out work with no delay at all, until the pot is full and five units of work have been completed. Every two 
seconds after that, a unit will drip out of the bucket and another unit of work can be completed.

Over the course of ten seconds, our example should process ten units of work; five initially and one every two seconds
for the remainder of the time.

## Example

```
func main() {

    // Create a new pot with a capacity of five, and which drips out a unit every two seconds
    pot := NewCreditsPot(CreditsPotConfig{ Size: 5, DripSeconds: 2 })

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