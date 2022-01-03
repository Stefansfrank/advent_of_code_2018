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

// input parser 
func parseFile (lines []string) (mp [][]byte) {

	mp = make([][]byte, len(lines)+2)
	mp[0] = make([]byte, len(lines[0])+2)
	mp[len(lines)+1] = make([]byte, len(lines[0])+2)
	for y, line := range lines {
		mp[y+1] = make([]byte, len(line)+2)
		for x, c := range line {
			mp[y+1][x+1] = byte(c)
		}
	}

	return
}

const emp = '.' // empty
const trs = '|' // trees
const yrd = '#' // lumber yard

// printing the map for debugging
func dump(mp [][]byte) {
	for y := 0; y < len(mp); y++ {
		for x := 0; x < len(mp[y]); x++ {
			fmt.Printf("%c",mp[y][x])
		}
		fmt.Print("\n")
	}
}

// counts all trees and lumber yards
func count(mp [][]byte) (trsCnt, yrdCnt int) {
	for y := 0; y < len(mp); y++ {
		for x := 0; x < len(mp[y]); x++ {
			if mp[y][x] == trs {
				trsCnt += 1
			} else if mp[y][x] == yrd {
				yrdCnt += 1
			}
		}
	}
	return
}

// one step in the simulation
func step(mp [][]byte) (nmp [][]byte) {

	// create padded new map
	nmp = make ([][]byte, len(mp))
	nmp[0] = make([]byte, len(mp[0]))
	nmp[len(mp)-1] = make([]byte, len(mp[0]))

	// go through map
	for y := 1; y < len(mp) - 1; y++ {
		nmp[y] = make([]byte, len(mp[y]))
		for x := 1; x < len(mp[y]) - 1; x++ {

			// go through adjacent and count
			var trCnt, ydCnt int
			for dy := -1; dy < 2; dy ++ {
				for dx := -1; dx < 2; dx ++ {
					if mp[y+dy][x+dx] == trs {
						trCnt += 1
					} else if mp[y+dy][x+dx] == yrd {
						ydCnt += 1
					}
				}
			}

			// core logic 
			switch mp[y][x] {
			case emp:
				if trCnt > 2 {
					nmp[y][x] = trs
				} else {
					nmp[y][x] = emp
				}
			case trs:
				if ydCnt > 2 {
					nmp[y][x] = yrd
				} else {
					nmp[y][x] = trs
				}
			case yrd:
				// note that I do count the tile itself not just adjacent 8
				if ydCnt == 1 || trCnt == 0 {
					nmp[y][x] = emp
				} else {
					nmp[y][x] = yrd
				}
			}
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
	input  := readTxtFile("d18." + dataset + ".txt")
	mp     := parseFile(input)

	// part 1
	for i := 0; i < 10; i++ {
		mp = step(mp)
	}
	trs, yrds := count(mp)
	fmt.Printf("Count product after 10 iterations: %v (%vx%v)\n", trs*yrds, trs, yrds)

	// part 2 - run up to 'tresh' iterations assuming a cycle has formed by then
	thresh := 1000
	for i := 10; i < thresh; i++ {
		mp = step(mp)
	}

	// detect cycle length and at the same time build a map of the results for the cycle
	trs, yrds = count(mp)
	sm1000   := trs*yrds
	cyc    := []int{ sm1000 }
	cycLen := 1
	for i := 0; true; i++ {
		mp = step(mp)
		trs, yrds = count(mp)
		if trs*yrds == sm1000 {
			cycLen = len(cyc)
			break
		} else {
			cyc = append(cyc, trs*yrds)
		}
	}

	// offset of the cycle value index to the modulo of the cycle length
	delta := thresh % cycLen
	fmt.Printf("And after 1000000000 iterations: %v\n", cyc[((1000000000 % cycLen) + cycLen - delta) % cycLen])

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}
