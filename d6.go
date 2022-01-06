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

// inline capable / no error handling Atoi
func atoi (x string) int {
	y, _ := strconv.Atoi(x)
	return y
}

// max / min
func min(i int, i2 int) int {
	if i < i2 {
		return i
	}
	return i2
}

// max / min
func max(i int, i2 int) int {
	if i < i2 {
		return i2
	}
	return i
}

// absolute
func abs(i int) int {
	if i > -1 {
		return i
	}
	return -i
}

// observation structure
type xy struct {
	x int
	y int
}

// input parser also returns the minimum and maximum dimensions
func parseFile (lines []string) (cnts []xy, mind, maxd xy) {

	cnts = make([]xy, len(lines))
	maxd = xy{0,0}
	mind = xy{1000,1000}

	for i, line := range lines {
		ix := strings.Index(line, ",")
		cnts[i] = xy{ atoi(line[:ix]), atoi(line[ix+2:]) }
		mind.x = min(mind.x, cnts[i].x)
		mind.y = min(mind.y, cnts[i].y)
		maxd.x = max(maxd.x, cnts[i].x)
		maxd.y = max(maxd.y, cnts[i].y)
	}
	return
}

// detects the relevant areas for part 1 and whether they are infinite
// also counts the total safe area for part 2
func detectArea(cnts []xy, mind, maxd xy, limit int) (area []int, inf []bool, safe int) {
	area = make([]int, len(cnts))
	inf  = make([]bool, len(cnts))
	var ix, mini, dist, totdist int

	for y := mind.y; y <= maxd.y; y ++ {
		for x := mind.x; x <= maxd.x; x ++ {

			// determines ix of closest center
			// and adds up the distance to all in the process
			ix = -1
			mini = maxd.x + maxd.y + 100
			totdist = 0
			for i, cnt := range cnts {
				dist = abs(cnt.y - y) + abs(cnt.x - x)
				if dist < mini {
					mini = dist
					ix   = i
				} else if dist == mini {
					ix = -1
				}
				totdist += dist
			}

			// if there is one closest center
			if ix > -1 {
				if x == mind.x || y == mind.y || x == maxd.x || y == maxd.y {
					inf[ix] = true
				}
				area[ix] += 1
			}

			// is the total distance under limit
			if totdist < limit {
				safe += 1
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

	var limit int
	if dataset[0:4] == "test" {
		limit = 32
	} else {
		limit = 10000
	}

	start  := time.Now()
	input  := readTxtFile("d6." + dataset + ".txt")

	cnts, mind, maxd := parseFile(input)
	area, inf, safe  := detectArea(cnts, mind, maxd, limit)

	maxArea := 0
	for i, a := range area {
		if a > maxArea && !inf[i] {
			maxArea = a
		}
	}

	fmt.Println("Max finite area: ", maxArea)
	fmt.Println("Max safe area: ", safe)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}