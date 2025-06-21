package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var counter int = 10
	numberChanel := make(chan int)
	squareChanel := make(chan int)

	wg.Add(2)

	go func() {
		defer wg.Done()

		numbers := make([]int, counter)
		for i := 0; i < counter; i++ {
			numbers[i] = rand.Intn(100)
		}

		for _, num := range numbers {
			numberChanel <- num
		}

		close(numberChanel)
	}()

	go func() {
		defer wg.Done()

		for num := range numberChanel {
			squareChanel <- num * num
		}

		close(squareChanel)
	}()

	for num := range squareChanel {
		fmt.Print(num, " ")
	}

	wg.Wait()
}
