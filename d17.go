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

func atoi(s string) (i int) {
	i, _ = strconv.Atoi(s)
	return
}

func min(i, j int) int {
	if i<j {
		return i
	}
	return j
}

func max(i, j int) int {
	if i>j {
		return i
	}
	return j
}

// input parser using Regex returns a map using the ints defined in constants below and the dimensions
func parseFile (lines []string) (mp [][]byte, xdim, ydim dim) {

	re  := regexp.MustCompile(`([x|y])=([0-9]+), [x|y]=([0-9]+)..([0-9]+)`)

	ls := make([][]int, len(lines))
	xdim = dim{500, 500}
	ydim = dim{12, 12}

	for i, line := range lines {

		match := re.FindStringSubmatch(line)
		ls[i] = make([]int, 4)
		ls[i][0] = int(match[1][0]) - int('x')
		for j := 1; j < 4; j++ {
			ls[i][j] = atoi(match[j+1])
		}

		if ls[i][0] == 0 {
			xdim.from = min(xdim.from, ls[i][1])
			xdim.to = max(xdim.to, ls[i][1])
			ydim.from = min(min(ydim.from, ls[i][2]), ls[i][3])
			ydim.to = max(max(ydim.to, ls[i][2]), ls[i][3])
		} else {
			ydim.from = min(ydim.from, ls[i][1])
			ydim.to = max(ydim.to, ls[i][1])
			xdim.from = min(min(xdim.from, ls[i][2]), ls[i][3])
			xdim.to = max(max(xdim.to, ls[i][2]), ls[i][3])
		}
	}
	xdim.from -= 1
	xdim.to   += 1

	mp = make([][]byte, ydim.to + 2)
	for y := 0; y < ydim.to + 2; y ++ {
		mp[y] = make([]byte, xdim.to - xdim.from + 1)
	}

	for _, l := range ls {
		if l[0] == 0 {
			for y := min(l[2],l[3]); y <= max(l[2], l[3]); y++ {
				mp[y][l[1]-xdim.from] = wll
			}
		} else {
			for x := min(l[2],l[3]); x <= max(l[2], l[3]); x++ {
				mp[l[1]][x-xdim.from] = wll
			}			
		}
	}

	// the well and first flow
	mp[0][500 - xdim.from] = src
	mp[1][500 - xdim.from] = flw

	return
}

type dim struct {
	from, to int
}

type xy struct {
	x, y int
}

const emp = byte(0) // empty
const flw = byte(1) // flowing water
const stw = byte(2) // standing water
const wll = byte(3) // wall
const src = byte(4) // source (well)

// printing the map for debugging
func dump(mp [][]byte) {
	for y := 0; y < len(mp); y++ {
		for x := 0; x < len(mp[y]); x++ {
			switch mp[y][x] {
			case emp:
				fmt.Print(".")
			case flw:
				fmt.Print("|")
			case stw:
				fmt.Print("~")
			case wll:
				fmt.Print("#")
			case src:
				fmt.Print("+")
			}
		}
		fmt.Print("\n")
	}
}

// counts all water or only standing water
func count(mp [][]byte, xdim, ydim dim, stonly bool) (cnt int) {
	for y := ydim.from; y <= ydim.to; y++ {
		for x := 0; x < len(mp[y]); x++ {
			if mp[y][x] == stw || (mp[y][x] == flw && !stonly) {
				cnt++
			}
		}
	}
	return
}

// the flow simulation (DFS path building)
func flow(mp [][]byte, xdim, ydim dim) {

	// the ist of active flows needing to be developed
	act := []xy{ xy{500 - xdim.from, 1} }

	// flow currently worked on
	cur := 0
	for cur < len(act) {

		// if there is empty below -> fall
		if mp[act[cur].y + 1][act[cur].x] == emp {
			if act[cur].y < ydim.to {
				act[cur].y += 1
				mp[act[cur].y][act[cur].x] = flw
			} else {
				cur += 1
				continue
			}

		// we hit flowing water (previously visited) -> this flow is finished
		} else if mp[act[cur].y + 1][act[cur].x] == flw {
			cur += 1
			continue

		// there is a wall / standing water below us
		} else {

			// detect horizontal extension of the barrier
			var xf, xt int
			ly     := act[cur].y
			lw, rw := false, false

			// how far to the left?
			for xf = act[cur].x; mp[ly+1][xf] > flw; xf-- {
				if mp[ly][xf] == wll {
					lw = true // hit a wall
					xf += 1
					break
				}
			}

			// how far to the right?
			for xt = act[cur].x; mp[ly+1][xt] > flw; xt++ {
				if mp[ly][xt] == wll {
					rw = true // hit a wall
					xt -= 1
					break
				}
			}

			// walls left and right -> standing water and set the flow one higher
			if lw && rw {
				for x := xf; x <= xt; x++ {
					mp[ly][x] = stw
				}
				act[cur].y -= 1

			// wall right -> floating water and this flow continues on left edge
			} else if rw {
				for x := xf; x <= xt; x++ {
					mp[ly][x] = flw
				}
				act[cur].x = xf

			// wall left -> floating water and this flow continues on right edge
			} else if lw {
				for x := xf ; x <= xt; x++ {
					mp[ly][x] = flw
				}
				act[cur].x = xt	

			// no walls -> floating water, this path continues on right edge and left edge is added to the flow list
			} else {
				for x := xf; x <= xt; x++ {
					mp[ly][x] = flw
				}
				act = append(act, xy{ xf, ly})
				act[cur].x = xt							
			}
		}
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
	input  := readTxtFile("d17." + dataset + ".txt")
	mp, xdim, ydim := parseFile(input)

	// both parts are solved in one simulatiion
	flow(mp, xdim, ydim)
	fmt.Println("Water reaches", count(mp, xdim, ydim, false), "spots")
	fmt.Println("Water standing in", count(mp, xdim, ydim, true), "spots")

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}