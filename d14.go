package main

import (
	"fmt"
	"strconv"
	"time"
)

// simulates n steps and prints the 10 digits after n
func sim(n int) {

	rcp := make([]byte, 2, n + 12) // I am using append but size it correctly from the start
	rcp[0] = byte(3)
	rcp[1] = byte(7)
	elf := []int{0,1}
	lmt := n + 10
	nw := byte(0)

	for len(rcp) < lmt  {
		nw = rcp[elf[0]] + rcp[elf[1]]
		if nw > byte(9) {
			rcp = append(rcp, byte(1))
			nw -= 10
		}
		rcp = append(rcp, nw)

		for i := 0; i < 2; i++ {
			elf[i] = (elf[i] + 1 + int(rcp[elf[i]])) % len(rcp)		
		}
	}

	fmt.Print("\nThe next ten digits after ", n, " are ")
	for i := 0; i < 10; i++ {
		fmt.Print(rcp[n+i])
	}
	fmt.Println()
}

// iterates until it matches and prints the amount of recipes before the match
func match(n string) {

	b1 := byte(1)
	b9 := byte(9)
	b10 := byte(10)
	rcp := make([]byte, 2)
	rcp[0] = byte(3)
	rcp[1] = byte(7)
	elf := []int{0,1}
	nw  := byte(0)

	// convert the pattern into []bytes
	tst := make([]byte, len(n))
	for i,c := range n {
		tst[i] = byte(c) - '0'
	}

	mtch := 0 // a running counter of the sequential digit matches
	var cnt int
	for cnt = 0; mtch < len(tst); cnt++ {
		nw = rcp[elf[0]] + rcp[elf[1]]
		if nw > b9 {
			rcp = append(rcp, b1)
			if tst[mtch] == b1 {
				mtch += 1
			} else {
				if mtch > 0 && tst[0] == b1 { // if an ongoin match was broken, maybe it matches the first
					mtch = 1
				} else {
					mtch = 0
				}
			}
			nw -= b10
		}
		if mtch < len(tst) {
			rcp = append(rcp, nw)
			if tst[mtch] == nw {
				mtch += 1
			} else {
				if mtch > 0 && tst[0] == nw { // if an ongoing match was broken, maybe it matches the first
					mtch = 1
				} else {
					mtch = 0
				}
			}
		}

		for i := 0; i < 2; i++ {
			elf[i] = (elf[i] + 1 + int(rcp[elf[i]])) % len(rcp)		
		}
	}

	fmt.Printf("%v recipes are to the left of the match\n\n", len(rcp) - len(tst))
}

func main() {
	
	start := time.Now()

	n    := "409551"
	ni,_ := strconv.Atoi(n)
	sim(ni)  // Part 1
	match(n) // Part 2

	fmt.Println("Execution time:", time.Since(start))
}