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

// runs the code with the instruction pointer logic
// 'zero' sets the initial value for r0
// max limits the amount of cpu steps no matter whether it terminates
func run(code []inst, ip int, zero int, max int) (reg []int) {

	reg = []int{zero,0,0,0,0,0}

	for cnt := 0; reg[ip] < len(code) && cnt < max; cnt++{
		if reg[ip] == 7 {
			fmt.Println(reg[5],"(",reg[2],")")
		}
		reg = cpu(code[reg[ip]], reg)
		reg[ip] += 1
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

	start  := time.Now()
	input  := readTxtFile("d19." + dataset + ".txt")
	code, ip := parseFile(input)
	
	// Part 1
	// Runs the code 25 steps long as register 4 is computed by then
	reg    := run(code, ip, 0, 25)

	// register 4 is target of the factor search
	facTgt := reg[4]

	// brute force the factor search
	reg0   := 0
	for i := 1; i <= facTgt; i++ {
		if facTgt % i == 0 {
			reg0 += i
		}
	}
	fmt.Println("\nRunning the code with register 0 = 0 leaves register 0 as", reg0)

	// Part 2 - identical solution to the above other than setting register 0 to 1
	reg     = run(code, ip, 1, 25)
	facTgt  = reg[4]
	reg0    = 0
	for i := 1; i <= facTgt; i++ {
		if facTgt % i == 0 {
			reg0 += i
		}
	}
	fmt.Println("\nRunning the code with register 0 = 1 leaves register 0 as", reg0, "\n")

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}