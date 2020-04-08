package main

import (
	"fmt"
	"strconv"
)

func main() {

	ch1 := make(chan int, 5)
	ch2 := make(chan string, 5)
	ch3 := make(chan []byte)

	ch1 <- 1
	ch1 <- 2
	ch1 <- 3
	close(ch1)

	go f1(ch1, ch2)
	go f2(ch2, ch3)
	f3(ch3)
}

func f1(in chan int, out chan string) {
	for i := range in {
		out <- strconv.Itoa(i)
	}
	close(out)
}

func f2(in chan string, out chan []byte) {
	for s := range in {
		out <- []byte(s)
	}
	close(out)
}

func f3(in chan []byte) {
	for b := range in {
		fmt.Println(string(b))
	}
}
