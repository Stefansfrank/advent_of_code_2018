package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	"strings"
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

// quick helper
func atoi (s string) (i int) {
	i, _ = strconv.Atoi(s)
	return
}

// input parser 
func parseFile (lines []string) (code []inst, ip int) {

	code  = []inst{}

	// the instruction pointer register
	ip = atoi(lines[0][4:])

	// the actual code
	for i := 1; i < len(lines); i++ {
		splt := strings.Split(lines[i], " ")
		ln   := inst{opc: splt[0], prm: make([]int, 3)}
		for j := 0; j <3; j++ {
			ln.prm[j] = atoi(strings.TrimSpace(splt[j+1]))
		}
		code = append(code, ln)		
	}

	return
}

// one instruction line for the cpu
type inst struct {
	opc string
	prm []int
}

const A = 0
const B = 1
const C = 2

// the cpu itself
func cpu(cdln inst, reg []int) (nreg []int) {
	nreg = make([]int, 6)
	copy(nreg, reg)

	switch cdln.opc {
	case "addr":
		nreg[cdln.prm[C]] = reg[cdln.prm[A]] + reg[cdln.prm[B]]
	case "addi":
		nreg[cdln.prm[C]] = reg[cdln.prm[A]] + cdln.prm[B]
	case "mulr":
		nreg[cdln.prm[C]] = reg[cdln.prm[A]] * reg[cdln.prm[B]]
	case "muli":
		nreg[cdln.prm[C]] = reg[cdln.prm[A]] * cdln.prm[B]
	case "banr":
		nreg[cdln.prm[C]] = reg[cdln.prm[A]] & reg[cdln.prm[B]]
	case "bani":
		nreg[cdln.prm[C]] = reg[cdln.prm[A]] & cdln.prm[B]
	case "borr":
		nreg[cdln.prm[C]] = reg[cdln.prm[A]] | reg[cdln.prm[B]]
	case "bori":
		nreg[cdln.prm[C]] = reg[cdln.prm[A]] | cdln.prm[B]
	case "setr":
		nreg[cdln.prm[C]] = reg[cdln.prm[A]] 
	case "seti":
		nreg[cdln.prm[C]] = cdln.prm[A]
	case "gtir":
		if cdln.prm[A] > reg[cdln.prm[B]] {
			nreg[cdln.prm[C]] = 1
		} else {
			nreg[cdln.prm[C]] = 0
		}
	case "gtri":
		if reg[cdln.prm[A]] > cdln.prm[B] {
			nreg[cdln.prm[C]] = 1
		} else {
			nreg[cdln.prm[C]] = 0
		}
	case "gtrr":
		if reg[cdln.prm[A]] > reg[cdln.prm[B]] {
			nreg[cdln.prm[C]] = 1
		} else {
			nreg[cdln.prm[C]] = 0
		}
	case "eqir":
		if cdln.prm[A] == reg[cdln.prm[B]] {
			nreg[cdln.prm[C]] = 1
		} else {
			nreg[cdln.prm[C]] = 0
		}
	case "eqri":
		if reg[cdln.prm[A]] == cdln.prm[B] {
			nreg[cdln.prm[C]] = 1
		} else {
			nreg[cdln.prm[C]] = 0
		}
	case "eqrr":
		if reg[cdln.prm[A]] == reg[cdln.prm[B]] {
			nreg[cdln.prm[C]] = 1
		} else {
			nreg[cdln.prm[C]] = 0
		}
	}
	return
}

// runs the code until it encounters line 'tstStp' 
// and prints the content of register 'tstRg' 
// 'zero' sets the initial value for r0
func runToTest(code []inst, ip int, tstRg int, tstStp int, zero int) (reg []int) {

	reg = []int{zero,0,0,0,0,0}

	for cnt := 0; reg[ip] < len(code); cnt++{
		if reg[ip] == tstStp {
			fmt.Println("Encountered test step, value of test register:", reg[tstRg])
			return
		}
		reg = cpu(code[reg[ip]], reg)
		reg[ip] += 1
	}

	return
}

// runs the code until it encounters the same value in register 'tstRg' than before in step 'tstStp'
// 'zero' sets the initial value for r0
func runToRepeat(code []inst, ip int, tstRg int, tstStp int, zero int) (reg []int) {

	tst  := map[int]bool{}
	last := 0

	reg = []int{zero,0,0,0,0,0}

	for cnt := 0; reg[ip] < len(code); cnt++{
		if reg[ip] == tstStp {
			if tst[reg[tstRg]] {
				fmt.Println("Encountered repeat - last test value before:", last)
				return
			} else {
				tst[reg[tstRg]] = true
				last = reg[tstRg]
			}
		}
		reg = cpu(code[reg[ip]], reg)
		reg[ip] += 1
	}

	return
}

// MAIN ----
// This code is very specialized for the test input at hand. The input code uses register 0 only to check a test
// register (in my case register 1) for equality and halts the program if they are equal. Thus the task at hand is
// observing the behavior of that test register since we can't impact anything other than providing the right test value.
// Part 1 - check the value of the test register the first time the test is encountered
// Part 2 - given the puzzle, the assumption is that eventually, a cyclic behavior is forming. 
// If that is the case, then the highest value that still halts is the last value before the first repeated value
func main () {

	dataset := ""
	if len(os.Args) < 2 || os.Args[1] == "" {
		fmt.Println("No argument given - trying 'test' dataset.")
		dataset = "test"
	} else {
		dataset = os.Args[1]
	}

	start  := time.Now()
	input  := readTxtFile("d21." + dataset + ".txt")
	code, ip := parseFile(input)
	
	// Part 1 - in my input, register 0 is tested in line 28 against register 1, thus the parameters 
	runToTest(code, ip, 1, 28, 0)

	// Part 2 - in my input, register 0 is tested in line 28 against register 1, thus the parameters 
	// runs long (2 minutes on my MP Pro !!)
	runToRepeat(code, ip, 1, 28, 0)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}