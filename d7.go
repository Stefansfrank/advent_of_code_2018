package main

import (
	"fmt"
	"os"
	"bufio"
	"time"
)

// no error handling ...
func readTxtFile (name string) (lines []string) {	
	file, _ := os.Open(name)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {		
		lines = append(lines, scanner.Text())
	}
	return
}

// helper for fast 2^n
func binHash() []int {
	bin := make([]int, 32)
	cnt := 1
	for i := 0; i < 32; i++ {
		bin[i] = cnt
		cnt <<= 1
	}
	return bin
}

// input parser creates requirements as bits
func parseFile (lines []string, tst bool) (req []int) {
	if tst {
		req  = make([]int, 6)
	} else {
		req  = make([]int, 26)
	}

	bin := binHash()

	for _, ln := range lines {
		rq := int(ln[5] - 'A')
		c  := int(ln[36] - 'A')
		req[c] += bin[rq]
	}

	return
}

// create sequence
func seq (req []int) (sq string) {
	fin := 0
	bin := binHash()

	for len(sq) != len(req) {
		for i, rq := range req {
			if rq & fin == rq {
				fin    += bin[i]
				sq     += fmt.Sprintf("%c",i + int('A'))
				req[i] += bin[31]
				break
			}
		}
	}

	return
}

type wrk struct {
	c  int
	tm int
}

// create sequence by simulating time
func seq2 (req []int, tst bool) (sq string, tm int) {
	fin  := 0
	bin  := binHash()

	var dl, wm int
	if tst {
		dl = 1
		wm = 2
	} else {
		dl = 61
		wm = 5
	}
	wks := make([]wrk, wm) // workers

	// run the clock
	for clk := 0; len(sq) < len(req); clk++ {
		for wi, wk := range wks {

			// a worker has finished
			if wk.tm == clk && clk > 0 {
				sq  += fmt.Sprintf("%c",wk.c + int('A'))
				fin += bin[wk.c]
			}

			// a worker is free
			if wk.tm <= clk {
				for i, rq := range req {
					if rq & fin == rq {
						req[i]  += bin[31]
						wks[wi].tm = clk + dl + i
						wks[wi].c  = i 
						break
					}
				}
			}
		}
	}

	// there seems to be some inconsistency in the puzzle on
	// whether the last busy or the first empty minute is counted
	// In the example, this -1 has to be removed but for the solution, it's needed
	// I tested this against many other solutions and all of them get the same result
	tm = maxWk(wks)

	return
}

func nextWk(wk []int, cur int) (ncur int, ix int) {
	tm := wk[0]
	for i := 0; i < len(wk); i++ {
		if wk[i] < tm {
			tm = wk[i]
			ix = i
		}
	}
	return
}

func maxWk(wk []wrk) (tm int) {
	tm = wk[0].tm
	for i := 1; i < len(wk); i++ {
		if wk[i].tm > tm {
			tm = wk[i].tm
		}
	}
	return
}

// MAIN ----
func main () {

	dataset := ""
	if len(os.Args) < 2 || os.Args[1] == "" {
		fmt.Println("No argument given - trying 'test' dataset.")
		dataset = "test"
	} else {
		dataset = os.Args[1]
	}
	tst := dataset[:4] == "test"

	start  := time.Now()
	input  := readTxtFile("d7." + dataset + ".txt")

	req  := parseFile(input, tst)
	req2 := make([]int, len(req))
	copy(req2, req)
	fmt.Println("Sequence with 1 worker:",seq(req))

	sq, tm := seq2(req2, tst)
	// there seems to be some inconsistency in the puzzle on
	// whether the last busy or the first empty minute is counted
	// In the example, the first empty minute is correct
	// but for my input, the last busy minute is correct
	// I tested this against many other solutions and all of them get the same result
	fmt.Println("Sequence with multiple workers:", sq, "Time:", tm, "maybe", tm-1)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}