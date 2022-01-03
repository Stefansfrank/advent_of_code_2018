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

// 
type sample struct {
	inp, out []int
	opc int
	prm []int
}

// input parser using Regex
func parseFile (lines []string) (smpls []sample, prog [][]int) {

	smpls = []sample{}
	prog  = [][]int{}
	for i := 0; i < len(lines); i++ {
		if len(lines[i]) > 7 && lines[i][:7] == "Before:" {
			smpl := sample{inp:make([]int, 4), out:make([]int, 4), prm:make([]int, 3)}
			splt := strings.Split(lines[i][9:len(lines[i])-1], ",")
			for j,p := range splt {
				smpl.inp[j] = atoi(strings.TrimSpace(p))
			}
			splt = strings.Split(lines[i+1], " ")
			smpl.opc = atoi(splt[0])
			for j := 1; j < 4; j++ {
				smpl.prm[j-1] = atoi(strings.TrimSpace(splt[j]))
			}
			splt = strings.Split(lines[i+2][9:len(lines[i+2])-1], ",")			
			for j,p := range splt {
				smpl.out[j] = atoi(strings.TrimSpace(p))
			}
			smpls = append(smpls, smpl)
			i += 3
		} else {
			if len(lines) > 6 {
				plin := make([]int, 4)
				splt := strings.Split(lines[i], " ")
				for j, p := range splt {
					plin[j] = atoi(strings.TrimSpace(p))
				}
				prog = append(prog, plin)		
			}
		}
	}
	return
}

func atoi (s string) (i int) {
	i, _ = strconv.Atoi(s)
	return
}

const addr = 0
const addi = 1
const mulr = 2
const muli = 3
const banr = 4
const bani = 5
const borr = 6
const bori = 7
const setr = 8
const seti = 9
const gtir = 10
const gtri = 11
const gtrr = 12
const eqir = 13
const eqri = 14
const eqrr = 15

const A = 0
const B = 1
const C = 2

func cpu(opc int, prm []int, reg []int) (nreg []int) {
	nreg = make([]int, 4)
	copy(nreg, reg)

	switch opc {
	case addr:
		nreg[prm[C]] = reg[prm[A]] + reg[prm[B]]
	case addi:
		nreg[prm[C]] = reg[prm[A]] + prm[B]
	case mulr:
		nreg[prm[C]] = reg[prm[A]] * reg[prm[B]]
	case muli:
		nreg[prm[C]] = reg[prm[A]] * prm[B]
	case banr:
		nreg[prm[C]] = reg[prm[A]] & reg[prm[B]]
	case bani:
		nreg[prm[C]] = reg[prm[A]] & prm[B]
	case borr:
		nreg[prm[C]] = reg[prm[A]] | reg[prm[B]]
	case bori:
		nreg[prm[C]] = reg[prm[A]] | prm[B]
	case setr:
		nreg[prm[C]] = reg[prm[A]] 
	case seti:
		nreg[prm[C]] = prm[A]
	case gtir:
		if prm[A] > reg[prm[B]] {
			nreg[prm[C]] = 1
		} else {
			nreg[prm[C]] = 0
		}
	case gtri:
		if reg[prm[A]] > prm[B] {
			nreg[prm[C]] = 1
		} else {
			nreg[prm[C]] = 0
		}
	case gtrr:
		if reg[prm[A]] > reg[prm[B]] {
			nreg[prm[C]] = 1
		} else {
			nreg[prm[C]] = 0
		}
	case eqir:
		if prm[A] == reg[prm[B]] {
			nreg[prm[C]] = 1
		} else {
			nreg[prm[C]] = 0
		}
	case eqri:
		if reg[prm[A]] == prm[B] {
			nreg[prm[C]] = 1
		} else {
			nreg[prm[C]] = 0
		}
	case eqrr:
		if reg[prm[A]] == reg[prm[B]] {
			nreg[prm[C]] = 1
		} else {
			nreg[prm[C]] = 0
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

	start  := time.Now()
	input  := readTxtFile("d16." + dataset + ".txt")
	smpls, prog := parseFile(input)
	
	// Part 1
	totCnt := 0
	topc := make([][][]int, 16)
	for _, smpl := range smpls {

		cnt  := 0
		if len(topc[smpl.opc]) == 0 {
			topc[smpl.opc] = [][]int{}
		}
		ttopc := []int{}
		for i := 0; i < 16; i++ {
			out := cpu(i, smpl.prm, smpl.inp)
			equ := 1
			for j,n := range out {
				if smpl.out[j] != n {
					equ = 0
					break
				}
			}
			cnt += equ
			if equ == 1 {
				ttopc = append(ttopc, i)
			}
		}
		if cnt > 2 {
			totCnt += 1
		} 
		topc[smpl.opc] = append(topc[smpl.opc], ttopc)
	}
	fmt.Println("\nSolution to part 1:", totCnt, "times could there be at least 3 opcodes per sample.\n")

	// Part 2
	// solution by looking at the output of this:
	for i, tp := range topc {
		fmt.Printf("Sample opcode %2v produces the same output as my numbers: ", i)
		otest := map[int]bool{}
		for j := 0; j < 16; j++ {
			otest[j] = true
		}
		for _, tt := range tp {
			test := map[int]bool{}
			for _,t := range tt {
				test[t] = true
			}
			for j,_ := range otest {
				if !test[j] {
					delete(otest, j)
				}
			}
		}
		for i,_ := range otest {
			fmt.Print(i," ")
		}
		fmt.Println()
	}

	// solution for opcode mapping from looking at the output above
	mp := []int{4,0,14,8,12,7,10,9,6,5,13,15,11,1,3,2}

	// run the cpu
	reg := []int{0,0,0,0}
	for _, code := range prog {
		reg = cpu(mp[code[0]],code[1:4], reg)
	}
	fmt.Println("\nOutput registers after running the input:", reg, "\n")

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}