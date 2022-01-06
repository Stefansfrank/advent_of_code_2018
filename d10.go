package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	"regexp"
	"math"
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

// int abs function
func abs(x int) int {
	if (x < 0) {
		return -x
	}
	return x
}

// input parser using Regex
func parseFile (lines []string) (stars []star) {

	re   := regexp.MustCompile(`position=< *(-?[0-9]+), *(-?[0-9]+)> velocity=< *(-?[0-9]+), *(-?[0-9]+)>`)
	stars = make([]star, len(lines)) 

	for i, line := range lines {
		match   := re.FindStringSubmatch(line)
		stars[i] = star{x:atoi(match[1]), y:atoi(match[2]), vx:atoi(match[3]), vy:atoi(match[4])}
	}
	return
}

// star structure
type star struct {
	x,y int
	vx, vy int
}

// move
func (s *star) mov(n int) {
	s.x += n * s.vx
	s.y += n * s.vy
}

// run until the area is expanding again
func sim(stars []star) {

	var xmin, ymin, xmax, ymax, narea int
	area := math.MaxInt
	
	for cnt := 0; true; cnt++{
		xmin = math.MaxInt
		ymin = math.MaxInt
		xmax = math.MinInt
		ymax = math.MinInt
		for i, _ := range stars {
			stars[i].mov(1)
			if stars[i].x < xmin {
				xmin = stars[i].x
			}
			if stars[i].y < ymin {
				ymin = stars[i].y
			}
			if stars[i].x > xmax {
				xmax = stars[i].x
			}
			if stars[i].y > ymax {
				ymax = stars[i].y
			}
		}
		narea = (xmax - xmin) * (ymax - ymin)

		// area is expanding
		if narea > area {

			// move all stars backwards 1
			for i, _ := range stars {
				stars[i].mov(-1)
			}

			// print result and exit
			dump(stars, xmin, xmax, ymin, ymax)
			fmt.Println("Time:", cnt)
			return

		// area is still contracting
		} else {
			area = narea
		}
	}
}

// prints the map of starts
func dump(sts []star, xmin, xmax, ymin, ymax int) {

	// create a map starting at 0,0
	xtop := xmax + 1 - xmin
	ytop := ymax + 1 - ymin
	mp := make([][]int, ytop)
	for y := 0; y < ytop; y++ {
		mp[y] = make([]int, xtop)
	}

	// add the stars with zero offset
	for _, s := range sts {
		mp[s.y - ymin][s.x - xmin] = 1
	}

	// print map
	for y := 0; y < len(mp); y++ {
		for x := 0; x < len(mp[0]); x ++ {
			fmt.Printf("%c", mp[y][x] * ('#' - '.') + '.')
		}
		fmt.Println()
	}
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
	input  := readTxtFile("d10." + dataset + ".txt")
	stars  := parseFile(input)

	sim(stars)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}