package main

import (
	"fmt"
	"os"
	"bufio"
	"time"
	"sort"
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

// track parts
const vline = 1
const hline = 2
const cross = 3
const ulcrv = 4
const urcrv = 5

// direction
const up  = 0
const rgt = 1
const dwn = 2
const lft = 3
var vdir []xy

type xy struct {
	x int
	y int
}

type car struct {
	pos xy
	dir int
	trn int
}

// parsing 
func parseFile(lines []string) (mp [][]int, cars map[xy]*car) {

	cars = map[xy]*car{}
	mp   = make([][]int, len(lines))
	var loc xy
	for y, line := range lines {
		mp[y] = make([]int, len(line))
		for x, c := range line {
			switch byte(c) {
			case '|':
				mp[y][x] = vline
			case '-':
				mp[y][x] = hline
			case '/':
				mp[y][x] = urcrv
			case '\\':
				mp[y][x] = ulcrv
			case '+':
				mp[y][x] = cross
			case '^':
				mp[y][x] = vline
				loc = xy{x,y}
				cars[loc] = &car{pos: loc, dir: up, trn:-1}
			case 'v':
				mp[y][x] = vline
				loc = xy{x,y}
				cars[loc] = &car{pos: loc, dir: dwn, trn:-1}
			case '<':
				mp[y][x] = hline
				loc = xy{x,y}
				cars[loc] = &car{pos: loc, dir: lft, trn:-1}
			case '>':
				mp[y][x] = hline
				loc = xy{x,y}
				cars[loc] = &car{pos: loc, dir: rgt, trn:-1}
			}
		}
	}
	return
}

// a single step for all cars
func step(mp [][]int, cars map[xy]*car) (crsh bool) {

	// determine the location of all cars
	locs := make([]xy, len(cars))
	i := 0
	for loc, _ := range cars {
		locs[i] = loc
		i++
	}

	// sort the cars in order to sequence their movement correctly
	sort.Slice(locs, func(i,j int) bool {
		if locs[i].y == locs[j].y {
			return locs[i].x < locs[j].x
		}
		return locs[i].y < locs[j].y
	})

	// loop though all cars
	for _, loc := range locs {
		cr, fnd := cars[loc]
	
		// car might have been removed due to crash
		if !fnd {
			continue
		}

		// move car
		cr.pos.x += vdir[cr.dir].x
		cr.pos.y += vdir[cr.dir].y 
		_,fnd = cars[cr.pos]

		// and determine new orientation
		switch mp[cr.pos.y][cr.pos.x] {
		case cross:
			cr.dir = (cr.dir + 4 + cr.trn) % 4
			cr.trn = (cr.trn + 2) % 3 - 1
		case ulcrv:
			if cr.dir == up || cr.dir == dwn {
				cr.dir = (cr.dir + 3) % 4
			} else {
				cr.dir = (cr.dir + 1) % 4
			}
		case urcrv:
			if cr.dir == up || cr.dir == dwn {
				cr.dir = (cr.dir + 1) % 4
			} else {
				cr.dir = (cr.dir + 3) % 4
			}
		}

		delete(cars, loc)
		if fnd { // collision
			delete(cars, cr.pos)
			fmt.Printf("Crash at (%v,%v)\n", cr.pos.x, cr.pos.y)
			crsh = true
		} else {
			cars[cr.pos] = cr
		}
	}
	return
}

// prints the map for debugging
func dump(mp [][]int, cars map[xy]*car) {

	trck := []string{" ","|","-","x","\\","/"}
	crs  := []string{"^",">","v","<"}
	for y,ln := range mp {
		for x, c := range ln {
			cr, fnd := cars[xy{x,y}]
			if fnd {
				fmt.Print(crs[cr.dir])
			} else {
				fmt.Print(trck[c])
			}
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

	start := time.Now()
	vdir   = []xy{ xy{x:0,y:-1}, xy{x:1,y:0}, xy{x:0,y:1}, xy{x:-1,y:0} }

	input  := readTxtFile("d13." + dataset + ".txt")
	mp, cars := parseFile(input)
	
	// part 1
	crsh := false
	for !crsh {
		crsh = step(mp, cars)
	}
	fmt.Println("---------------------------")

	// part 2
	for len(cars) > 1 {
		step(mp, cars)
	}
	for _, cr := range cars {
		fmt.Printf("---------------------------\nLast car at (%v,%v)\n", cr.pos.x, cr.pos.y)
	}

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}