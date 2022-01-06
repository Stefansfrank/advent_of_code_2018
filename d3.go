package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	"regexp"
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

// inline capable / no error handling Atoi
func atoi (x string) int {
	y, _ := strconv.Atoi(x)
	return y
}

// claim structure
type claim struct {
	no  int
	x,y int
	w,h int
}

// input parser using Regex
func parseFile (lines []string) (claims []claim) {

	re    := regexp.MustCompile(`#([0-9]+) @ ([0-9]+),([0-9]+): ([0-9]+)x([0-9]+)`)
	claims = make([]claim, len(lines))

	for i, line := range lines {
		match := re.FindStringSubmatch(line)
		claims[i] = claim{no:atoi(match[1]),
						   x:atoi(match[2]), y:atoi(match[3]),
						   w:atoi(match[4]), h:atoi(match[5])}	
	}

	return
}

const mpSize = 1000

// empty map to count overlaps
func newMap() [][]int {
	mp := make([][]int, mpSize)
	for i := 0; i < mpSize; i++ {
		mp[i] = make([]int, mpSize)
	}
	return mp
}

// the core logic of counting stuff
func countClaims(cls []claim) (int, int) {

	// counts the claims per field of the fabric
	mp := newMap()
	for _, cl := range cls {
		for y := cl.y; y < cl.y + cl.h; y++ {
			for x := cl.x; x < cl.x + cl.w; x++ {
				mp[y][x] += 1
			}
		}
	}

	// counts the fields with 2 or more claims
	cnt := 0
	for y := 0; y < mpSize; y++ {
		for x := 0; x < mpSize; x++ {
			if mp[y][x] > 1 {
				cnt += 1
			}
		}
	}

	// finds the claim that has no overlap
	cno := 0 
	clm: for _, cl := range cls {
		for y := cl.y; y < cl.y + cl.h; y++ {
			for x := cl.x; x < cl.x + cl.w; x++ {
				if mp[y][x] > 1 {
					continue clm
				}
			}
		}
		cno = cl.no 
		break
	}

	return cnt, cno
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
	input  := readTxtFile("d3." + dataset + ".txt")
	claims := parseFile(input)

	cnt, cno := countClaims(claims)
	fmt.Println("Overlap fields:", cnt, "\nClaim without overlap:", cno)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}