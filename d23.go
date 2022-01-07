package main

import (
	"fmt"
	"os"
	"bufio"
	"time"
	"regexp"
	"strconv"
	"sort"
)

// no error handling ...
func readFileStr (name string) (lines []string) {	
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
	i,_ = strconv.Atoi(s)
	return
}

// int abs function
func abs(x int) int {
	if (x < 0) {
		return -x
	}
	return x
}

// int min function
func min(i, j int) int {
    if i<j {
        return i
    }
    return j
}

// int max function
func max(i, j int) int {
    if i>j {
        return i
    }
    return j
}

// 3D objects --------------------------------------
type xyz struct {
    x, y, z int
}

// add
func (p xyz) add(v xyz) xyz {
	return xyz{p.x + v.x, p.y + v.y, p.z + v.z}
}

// sub
func (p xyz) sub(v xyz) xyz {
	return xyz{p.x - v.x, p.y - v.y, p.z - v.z}
}

// mult
func (p xyz) mult(m int) xyz {
	return xyz{p.x * m, p.y * m, p.z * m}
}

// equality
func (p xyz) equ(p2 xyz) bool {
	return p.x == p2.x && p.y == p2.y && p.z == p2.z
}

// pos
func (p xyz) pos() bool {
	return p.x > -1 && p.y > -1 && p.z > -1
}

// distance
func (p xyz) dist(p2 xyz) int {
	return abs(p2.x - p.x) + abs(p2.y - p.y) + abs(p2.z - p.z)
}

// range structure
type rng struct {
	fr, to int
}

// 3D box structure
type box struct {
	x, y, z rng
}

// ----------------------------------------------------

// create empty [][][]int array
func empty3dInt(dim xyz) (mp [][][]int) {
	mp = make([][][]int, dim.z)
	for z := 0; z < dim.z; z++ {
		mp[z] = make([][]int, dim.y)
		for y := 0; y < dim.y; y++ {
			mp[z][y] = make([]int, dim.x)
		}
	}
	return
}

// create empty [][][]byte array
func empty3dByte(dim xyz) (mp [][][]byte) {
	mp = make([][][]byte, dim.z)
	for z := 0; z < dim.z; z++ {
		mp[z] = make([][]byte, dim.y)
		for y := 0; y < dim.y; y++ {
			mp[z][y] = make([]byte, dim.x)
		}
	}
	return
}

// create empty [][][]bool array
func empty3dBool(dim xyz) (mp [][][]bool) {
	mp = make([][][]bool, dim.z)
	for z := 0; z < dim.z; z++ {
		mp[z] = make([][]bool, dim.y)
		for y := 0; y < dim.y; y++ {
			mp[z][y] = make([]bool, dim.x)
		}
	}
	return
}

// create empty [][][]string array
func empty3dString(dim xyz) (mp [][][]string) {
	mp = make([][][]string, dim.z)
	for z := 0; z < dim.z; z++ {
		mp[z] = make([][]string, dim.y)
		for y := 0; y < dim.y; y++ {
			mp[z][y] = make([]string, dim.x)
		}
	}
	return
}

// -----------------------------------------------------------------

// nanobot representation
type nanobot struct {
	loc xyz
	rad int
}

// input parser using Regex
// pos=<-33594389,69103993,46909087>, r=93546878
func parseFile (lines []string) (nbots []nanobot) {

	re  := regexp.MustCompile(`pos=<(-?[0-9]+),(-?[0-9]+),(-?[0-9]+)>, r=(-?[0-9]+)`)
	nbots = []nanobot{}
	for _, line := range lines {
		match := re.FindStringSubmatch(line)
		nbot := nanobot{ xyz{ atoi(match[1]), atoi(match[2]), atoi(match[3])}, atoi(match[4]) }
		nbots = append(nbots, nbot)
	}
	return
}

// Part 1 solution and determination of dims
func cntInRange(nbots []nanobot) (sm int, bx box) {

	sort.Slice(nbots, func (i, j int) bool { return nbots[i].rad > nbots[j].rad } )

	for _, nb := range nbots {
		if nb.loc.dist(nbots[0].loc) <= nbots[0].rad {
			sm += 1
			bx.x.fr = min(nb.loc.x, bx.x.fr)
			bx.y.fr = min(nb.loc.y, bx.y.fr)
			bx.z.fr = min(nb.loc.z, bx.z.fr)
			bx.x.to = max(nb.loc.x, bx.x.to)
			bx.y.to = max(nb.loc.y, bx.y.to)
			bx.z.to = max(nb.loc.z, bx.z.to)
		}
	}
	return
}

// divide all coorindates of a point by 'm'
func (p xyz) shrink(m int) xyz {
	return xyz{ p.x/m, p.y/m, p.z/m }
}

// finds out for point 'pt' by how many bots it is covered
// supports a divisor by which all coordinates other than the given point are divided by
func botsInRange(nbots []nanobot, pt xyz, div int) (cnt int) {
	for _, nb := range nbots {
		if pt.dist(nb.loc.shrink(div)) <= nb.rad/div {
			cnt++
		}
	}
	return
}

// finds best point by changing the granularity of the coordinate system
// first, all coordinates are divided by 10M which puts the bots into [-15..15]
// then I scane that cube of 31x31x31 for the best location. Then I increase the
// resolution by a factor of 10 and look into the 31x31x31 cube around the best solution
// in the coarser resolution. This is repeated until the divisor is 1
// it returns a slice of best points sorted with the one closest to 0,0,0 in [0]
func findBest(nbots []nanobot, dims box) (pts []xyz, max int) {

	div := 10000000   // initial divisor
	cnt := xyz{0,0,0} // initial center of the search area

	var sm int
	for div > 0 {

		fmt.Print(".")
		pts = []xyz{}
		max = 0
		for z := -15; z <= 15; z++ {
			for y := -15; y <= 15; y++ {
				for x := -15; x <= 15; x++ {
					sm = botsInRange(nbots, cnt.add(xyz{x,y,z}), div)
					if sm > max {
						pts = []xyz{ cnt.add(xyz{x,y,z}) }
						max = sm
					} else if sm == max {
						pts = append(pts, cnt.add(xyz{x,y,z}) )
					}
				}
			}
		}

		// sorts by proximity to 0,0,0
		sort.Slice(pts, func (i, j int) bool { return pts[i].dist(xyz{0,0,0}) < pts[j].dist(xyz{0,0,0}) })
		cnt = pts[0].mult(10)
		div /= 10
	}
	fmt.Println()

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
	input  := readFileStr("d23." + dataset + ".txt")
	nbots  := parseFile(input)

	cnt, dims := cntInRange(nbots)
	fmt.Println("\nDimensions:", dims, "\nSum of bots in range of largest:", cnt)
	pts, mx   := findBest(nbots, dims)
	fmt.Printf("One of the closest points with max coverage from %v bots is %v\nIt's distance from 0,0,0 is :%v\n\n", mx, pts[0], pts[0].dist(xyz{0,0,0}))

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}