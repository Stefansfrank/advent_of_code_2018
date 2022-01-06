package main

import (
	"fmt"
	"time"
)

// int max function
func max(i, j int) int {
    if i>j {
        return i
    }
    return j
}

// 2D objects --------------------------------------
type xy struct {
    x, y int
}

// add
func (p xy) add(v xy) xy {
	return xy{p.x + v.x, p.y + v.y}
}

// equality
func (p xy) equ(p2 xy) bool {
	return p.x == p2.x && p.y == p2.y
}

// pos
func (p xy) pos() bool {
	return p.x > -1 && p.y > -1
}

// create empty [][]int array
func empty2dInt(dim xy) (mp [][]int) {
	mp = make([][]int, dim.y)
	for i := 0; i < dim.y; i++ {
		mp[i] = make([]int, dim.x)
	}
	return
}

// create empty [][]bool array
func empty2dBool(dim xy) (mp [][]bool) {
	mp = make([][]bool, dim.y)
	for i := 0; i < dim.y; i++ {
		mp[i] = make([]bool, dim.x)
	}
	return
}

// puzzle specific code here --------------------------

// create erosion map 
func buildMap(dim, tgt xy, depth int) (mp [][]int) {
	mp = empty2dInt(dim)

	for i := 1; i < max(dim.x, dim.y); i++ {
		if i < dim.x {
			mp[0][i] = (i * 16807 + depth) % 20183
		}
		if i < dim.y {
			mp[i][0] = (i * 48271 + depth) % 20183
		}
	}

	mp[0][0] = (0 + depth) % 20183
	for y := 1; y < len(mp); y++ {
		for x := 1; x < len(mp[0]); x++ {
			if !tgt.equ(xy{x,y}) {
				mp[y][x] = (mp[y-1][x] * mp[y][x-1] + depth) % 20183
			} else {
				mp[y][x] = (0 + depth) % 20183
			}
		}
	}
	return
}

// dump map (for debugging)
func dumpMap(mp [][]int) {
	for _,ln := range mp {
		for _,c := range ln {
			switch c % 3 {
			case rock:
				fmt.Print(".")
			case wetl:
				fmt.Print("=")
			case nrrw:
				fmt.Print("|")
			}	
		}
		fmt.Println()
	}
}

// coding of terrain and equipment
const empty = 0
const torch = 1
const climb = 2
const rock = 0
const wetl = 1
const nrrw = 2

// calculates the risk (part 1)
func risk(mp [][]int, tgt xy) (sm int) {
	for y := 0; y <= tgt.y; y++ {
		for x := 0; x <= tgt.x; x++ {
			sm += mp[y][x] % 3
		}
	}
	return 
}

// the path structure we are using to compute and compare paths
type path struct {
	loc  xy // current locations
	cst int // cost (minutes)
	equ int // the currently equiped inventory
}

// in both the visited map and the open path map, we can not just index with x,y
// since a path that ends at a certain place with different equipment could be
// better or worse so we need to consider both. 
// Thus both the visited and open path maps are 3 dimensional using x,y,equipment
type locequ struct {
	loc xy 
	equ int
}

// dijkstra search of the best path
func dijk(mp [][]int, tgt xy) path {

	// the rules of equpiment usage
	allowed := empty2dBool(xy{3,3})
	allowed[rock][empty] = false
	allowed[rock][torch] = true
	allowed[rock][climb] = true
	allowed[wetl][empty] = true
	allowed[wetl][torch] = false
	allowed[wetl][climb] = true
	allowed[nrrw][empty] = true
	allowed[nrrw][torch] = true
	allowed[nrrw][climb] = false

	// potential neighbors
	vdir := []xy{ xy{0,-1}, xy{1,0}, xy{0,1}, xy{-1,0}}

	// temp variables used in the main loop
	var pth path
	var np  xy
	var npl locequ
	var nt, ncst int

	// the two core data structures for the search - the places visited and the open paths
	// 'pths' (open paths) contains only potential paths that have already identified 
	// as a neighbor of a visited location
	vis   := map[locequ]bool {}
	pths  := map[locequ]*path { locequ{ xy{0,0}, torch }: &path{ loc:xy{0,0}, cst:0, equ:torch } }

	// main loop
	for len(pths) > 0 {

		// find the lowest cost path, remove it from the open paths list and mark visited
		minCst := 10000000
		var minIx locequ
		for i, p := range pths {
			if p.cst < minCst {
				pth = *p
				minCst = p.cst
				minIx  = i
			}
		}
		delete(pths, minIx)
		vis[minIx] = true

		// actual exit condition - target is visited with a torch
		if minIx.loc.equ(tgt) && minIx.equ == torch {
			return pth
		}

		// the next options for the current visited path are the four neighbouring directions
		// and two potenital equipment changes (one will later be pruned out as not allowed)
		// from a graph perspective an equipment change is just one of the possible next steps
		for i := 0; i < 6; i++ {

			// new location
			if i < 4 {
				np  = pth.loc.add(vdir[i])
			 } else {
			 	np  = pth.loc
			 }

			// keep it inside mapped region
		 	if !np.pos() || np.x >= len(mp[0]) || np.y >= len(mp) {
		 		continue
		 	}

		 	// determine terrain at the new location
		 	nt   = mp[np.y][np.x] % 3

		 	// set the equipment and use it in order
		 	// to build the index for the new location
		 	npl  = locequ{np, pth.equ}
		 	ncst = 1
		 	if i > 3 {
		 		npl.equ = (npl.equ + i - 3) % 3
		 		ncst = 7
		 	}

		 	// don't consider if already visited or not allowed
		 	if vis[npl] || !allowed[nt][npl.equ] {
		 		continue
		 	}

		 	// already in the open path list?
		 	_, fnd := pths[npl]
		 	if fnd {

		 		// update cost if lower
		 		if pths[npl].cst > ncst + pth.cst {
		 			pths[npl].cst = ncst + pth.cst
		 		}

		 	} else {

				// create a new path and add to open paths
				npth := path{ loc:np, cst:pth.cst + ncst, equ:npl.equ }
				pths[npl] = &npth
			}
		} 
	}

	// returns null path if target is never reached
	return path{}
}

// MAIN ----
func main () {

	start  := time.Now()

	target := []xy{ xy{13, 743}, xy{10, 10} } // data given (input & example)
	depth  := []int{ 8112, 510 }              // data given (input & example)
	buffer := 100 // the buffer size beyond the target bound rectangle
	ix     := 0   // 0 = my input, 1 = example input

	mp := buildMap(target[ix].add(xy{buffer, buffer}), target[ix], depth[ix])

	// Part 1
	fmt.Println("Risk level:", risk(mp, target[ix]))

	// Part 2
	pth := dijk(mp, target[ix])
	fmt.Printf("Time to rescue: %v mins (%v:%v hrs)\n", pth.cst, pth.cst/60, pth.cst%60)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}