package main

import (
	"fmt"
	"time"
)

// Part 1: computes the grid of cells and their power
// also computes a grid with the sum of the 3x3 field with x,y as the top left corner
func compMap(sn int) (mp, mp33 [][]int) {

	var rid, pl, xx, yy, dx, dy int
	mp   = make([][]int, 301) // the grid with the values
	mp33 = make([][]int, 301) // the grid with 3x3 sums for each top left grid point
	for y := 1; y <= 300; y ++ {
		mp[y]   = make([]int, 301)
		mp33[y] = make([]int, 301)
		for x := 1; x <= 300; x ++ {

			// content
			rid      = x + 10
			pl       = (rid * y + sn) * rid
			mp[y][x] = ((pl % 1000) - (pl % 100)) / 100 - 5
			
			// add to 3x3 sums this point belongs to
			for dy = 0; dy < 3; dy++ {
				for dx = 0; dx < 3; dx++ {
					xx = x - dx
					yy = y - dy
					if xx > 0 && yy > 0 {
						mp33[yy][xx] += mp[y][x]
					}
				}
			}
		}
	} 
	fmt.Println()
	return
}

// used in part 1 to find the maximum of the 3x3 sum grid
func maxMap33(mp33 [][]int) (mx, xx, yy int) {
	for y := 1; y <= 298; y ++ {
		for x := 1; x <= 298; x ++ {
			if mp33[y][x] > mx {
				mx = mp33[y][x]
				xx = x
				yy = y 
			}
		}
	}
	return
}

// solves part 2 using 2d prefix sums to reduce order
func findMax(sn int) (max, xmax, ymax, szmax int) {

	var rid, pl int
	mp := make([][]int, 301)

	// calculate base grid itself 
	for y := 0; y <= 300; y++ {
		mp[y] = make([]int, 301)
		for x := 0; x <= 300; x++ {

			// content of cell itself 
			rid      = x + 10
			pl       = (rid * y + sn) * rid
			mp[y][x] = ((pl % 1000) - (pl % 100)) / 100 - 5
		}
	}

	// calculate the 2D prefix sums for each cell
	for y := 1; y <= 300; y ++ {
		for x := 1; x <= 300; x ++ {
			mp[y][x] = mp[y][x] + mp[y-1][x] + mp[y][x-1] - mp[y-1][x-1]
		}
	}
		
	// loop through map and compute all sizes using the prefix sums making it O(n^3)
	var sm int
	for sz := 1; sz <= 300; sz++ {
		for y := 1; y <= 300 - sz; y++ {
			for x := 1; x <= 300 - sz; x++ {
				sm = mp[y + sz][x + sz] - mp[y][x + sz] - mp[y + sz][x] + mp[y][x]
				if sm > max {
					max   = sm 
					xmax  = x + 1 // the box we computed is next to the top left value we used
					ymax  = y + 1 // the box we computed is next to the top left value we used
					szmax = sz
				} 
			}
		}
	} 

	return
}

// prints map for debugging
func dumpMap(mp [][]int, from,to int) {
	for y := from; y <= to; y ++ {
		for x := from; x <= to; x ++ {
			fmt.Printf("%3v", mp[y][x])
		}
		fmt.Println()
	}
}

func main() {
	
	start := time.Now()

	serial := 3463 // serial number (my input)

	_, mp33 := compMap(serial)
	_, x, y := maxMap33(mp33)
	fmt.Printf("Part 1: Coordinates of max 3x3: %v,%v\n", x, y)

	_, x, y, sz := findMax(serial)
	fmt.Printf("Part 2: Coordinates and size of max box: %v,%v,%v\n", x, y, sz)

	fmt.Println("Execution time:", time.Since(start))
}