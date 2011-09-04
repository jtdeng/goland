package main

import (
	"fmt"
)
// Copy the values from channel 'in' to channel 'out',
// removing those divisible by 'prime'.
func filter(in, out chan int, prime int) {
	for {
		i := <-in // Receive value of new variable 'i' from 'in'.
		if i%prime != 0 {
			out <- i // Send 'i' to channel 'out'.
		}
	}
}

func generate(ch chan int) {
	for i := 2; ; i++ {
		ch <- i // Send 'i' to channel 'ch'.
	}
}

/*
   func sieve() chan int {
       out := make(chan int)
       go func() {
           ch := generate(out)
           for {
               prime := <-ch
               out <- prime
               ch = filter(ch, prime)
           }
       }()
       return out
   }
*/
func main() {
	ch := make(chan int)       // Create a new channel.
	go generate(ch)            // Start generate() as a goroutine.
	for i := 0; i < 1000; i++ { // Print the first hundred primes.
		prime := <-ch
		fmt.Println(prime)
		ch1 := make(chan int)
		go filter(ch, ch1, prime)
		ch = ch1
	}

}
