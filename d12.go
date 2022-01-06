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

// typing helpers when dealing with bytes
const b0 = byte(0)
const b1 = byte(1)
const b2 = byte(2)
const b4 = byte(4)
const b31 = byte(31)

// parsing the input so that the pots are represented as 1s and 0s in a []byte
// and the rules are 5 bit ints (0-31) where a set bit expresses a # and a zero bit a .
// only rules that creata a plant are parsed since the default per pot is 'no plant'
func parseFile(lines []string) (pots, rules []byte) {

	pots = make([]byte, len(lines[0]) - 15)
	for i := 15; i < len(lines[0]); i++ {
		pots[i-15] = ('.' - lines[0][i])/('.' - '#')
	}

	rules = []byte{}
	for i := 2; i < len(lines); i++ {
		if lines[i][9] == '#' {
			pat := b0
			for j := 0; j < 5; j++ {
				pat <<= 1
				pat += ('.' - lines[i][j])/('.' - '#')
			}
			rules = append(rules, pat)
		}
	}

	return
}

// the simulation starting with configuration 'pots' and a set of rules and running 'n' generations
// it returns the state of all pots, the count at the end and the count of the generation before
func sim(pots, rules []byte, n int) (npots []byte, this, last int) {

	delta := (n+2) * 2
	inlen := len(pots)
	pots   = append(make([]byte, delta), pots...)
	pots   = append(pots, make([]byte, delta)...)

	sub   := b0

	var from, to, ix int

	// generations loop
	for it := 0; it < n; it++ {

		npots  = make([]byte, len(pots))
		from   = delta - 2*(it + 1)
		to     = inlen + delta + 2*(it + 1)

		// loop through potential pots
		for ix = from; ix <= to; ix++ {

			// the state of the relevant 5 pots as 5 bit int
			// (moves along by rotating bits, dropping bits 6+ and adding one new pot)
			sub = (sub << 1) & b31 + pots[ix + 2]

			// check whether a rule applies
			for _,r := range rules {
				if sub == r {
					npots[ix] = 1
					break
				}
			}
		}

		// for next iteration
		pots = npots

		// count the second to last iteration
		if it == n - 2 {
			last = count(pots, delta)
		}
	}
	this = count(pots, delta)
	return
}

// counts the pot sum assuming pot 'delta' to be zero
func count(pots []byte, delta int) (sm int) {
	for i,p := range pots {
		sm += int(p) * (i - delta)
	}
	return sm
}

// prints the pots for debugging
func dump(bs []byte) {
	std := false
	for _, b := range bs {
		if b == 1 {
			std = true
			fmt.Print("#")
		} else if std {
			fmt.Print(".")
		}
	}
	fmt.Println()
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
	input  := readTxtFile("d12." + dataset + ".txt")
	pots, rules := parseFile(input)

	_, this, _ := sim(pots, rules, 20)
	fmt.Println("Count after 20:", this)

	// I observed that the evolution of plants becomes linear after 1000 iterations
	// so calculating step 50000000000 is a simple linear equation extrapolating
	// from the differenc between the 999th and the 1000th iteration
	_, this, last := sim(pots, rules, 1000)
	mul := this - last
	add := this - 1000*mul
	fmt.Println("Count after 50000000000:", 50000000000*mul + add)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}