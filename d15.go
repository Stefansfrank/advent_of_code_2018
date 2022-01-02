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

// input parser 
func parseFile (lines []string) (mp [][]byte, ppl []*people, pploc map[xy]*people) {

	mp  = make([][]byte, len(lines))
	ppl = []*people{}
	pploc = map[xy]*people{}
	for y, line := range lines {
		mp[y] = make ([]byte, len(line))
		for x, c := range line {
			mp[y][x] = byte(c)
			if c == elf {
				pp := &people{true, 200, xy{x,y}}
				ppl = append(ppl, pp)
				pploc[xy{x,y}] = pp
			} else if c == gob {
				pp := &people{false, 200, xy{x,y}}
				ppl = append(ppl, pp)				
				pploc[xy{x,y}] = pp
			}
		}
	}

	return
}

// helper constants
const up  = 0
const lft = 1
const rgt = 2
const dwn = 3
var vdir []xy

const elf = 'E'
const gob = 'G'
const emp = '.'
const wal = '#'

// coordinate struct
type xy struct {
	x, y int
}

// all elfs and goblins
type people struct {
	elf bool
	hp  int
	loc xy
}

// a path that is evaluated for each elf/goblin
type path struct {
	sdir int  // start direction (def. see constants)
	len  int  // length of the path
	loc  xy   // target loc of the path
	dead bool // indicates whether the path is dead (should not be continued)
	vld  bool // indicates whether the path is valid i.e. leads to an enemy
}

// print out map for debugging
func dump(mp [][]byte) {
	for _,ln := range mp {
		for _,c := range ln {
			fmt.Printf("%c", c)
		}
		fmt.Println()
	}
	fmt.Println()
}

// returns map content and new coordinates for a neighbouring location in direction 'dir'
func neigh(loc xy, dir int, mp [][]byte) (bt byte, nloc xy) {
	nloc = xy{x:loc.x + vdir[dir].x, y:loc.y + vdir[dir].y}
	bt   = mp[nloc.y][nloc.x]
	return
}

// counts elfs in the people list
func countElfs(ppl []*people) (elfCnt int) {
	for _, pp := range ppl {
		if pp.elf {
			elfCnt++
		}
	}
	return
}

// this is the main logic ... and not pretty :()
// I am not sure the problem is solvable without some structure like this
// even though it could be better structured using subroutines
func battle(strength int, mp [][]byte, ppl []*people, pploc map[xy]*people) (elfCnt, totHp, cnt int) {

	// for speed reasons, I declare a lot of temporary
	// variables here so I save malloc time. The main structures are:
	var attack *people  // holds the enemy to attack as determined during path exploration
	var pths []path     // the structure collecint all possible paths a person could go
	var vis map[xy]int  // the visited map contains the index of the currently best path to location xy
	var nloc, tloc xy
	var npix int
	var more, fnd bool
	var tdir, thp, nlen int
	var died bool
	var nb byte

	// main tic / iterations
	for cnt = 0; true || cnt < 12; cnt++ {

		// sort the fighters in the correct order as they might have moved
		sort.Slice(ppl, func (i,j int) bool {
			if ppl[i].loc.y == ppl[j].loc.y {
				return ppl[i].loc.x < ppl[j].loc.x
			}
			return ppl[i].loc.y < ppl[j].loc.y })
		died = false

		// go through all of them
		for ppix, pp := range ppl {
			if pp.hp < 1 {
				continue
			}

			attack = nil
			thp    = 400
			pths = []path{}
			vis  = map[xy]int{}

			// determine up to 4 potential start paths
			// also detecting a potential immediate attack target
			// choosing the one with the lowest HP if there are multiple
			for i := 0; i < 4; i++ {
				nb, nloc = neigh(pp.loc, i, mp)
				if (nb == elf && !pp.elf) ||
					(nb == gob && pp.elf) {
						if pploc[nloc].hp < thp && pploc[nloc].hp > 0 {
							attack = pploc[nloc]
							thp = pploc[nloc].hp
						}
				} else if nb == emp {
					vis[nloc] = len(pths)
					pths = append(pths, path{sdir:i, len:1, loc:nloc})
				}
			}

			// if no imminent attack is possible, let's do a BFS identifying the shortest path to an enemy
			if attack == nil {

				// loops through the path finding until there is a valid path
				more = true
				for cnt2 := 0; more; cnt2++{

					// loop through valid paths adding possible continuations
					for pix, pth := range pths {

						if pth.vld || pth.dead {
							continue
						}

						// look at the four neighbours of the path
						for i := 0; i < 4; i++ {

							// is that neighbor open or an enemy?
							nb, nloc = neigh(pth.loc, i, mp)
							if nb == emp || (nb == gob && pp.elf) || (nb == elf && !pp.elf) {

								// is that neighbor visited and if so, is my path better (shorter or better start direction)
								npix, fnd = vis[nloc]
								nlen = pth.len + 1
								if !fnd || nlen < pths[npix].len || 
									(nlen == pths[npix].len && pth.sdir < pths[npix].sdir) {

									vis[nloc] = len(pths)	
									pths = append(pths, path{sdir:pth.sdir, len:nlen, loc:nloc})

									// if I am squeezing out an inferior solution, mark it dead so it's not continued
									if fnd {
										pths[npix].dead = true
										pths[npix].vld  = false
									}

									// if that was in fact an enemy mark this valid
									if (nb == gob && pp.elf) || (nb == elf && !pp.elf) {
										pths[len(pths)-1].vld = true

										// if I encountered the enemy after just one move, still attack
										if pth.len == 1 && pploc[nloc].hp < thp && pploc[nloc].hp > 0{
											attack = pploc[nloc]
											thp = attack.hp
										}
									}
								}
							}					
						}

						// I created all new paths from this path so it's no longer needed
						pths[pix].dead = true
					}

					// stop searching for path if there is at least one valid or only dead paths
					more = false 
					for _, pth := range pths {
						if pth.vld {
							more = false
							break
						}
						if !pth.dead {
							more = true
						}
					}
				}

				// now actually move that valid path
				// look if there are many valids and find the best
				// they would all be the same length ...
				// the best is the first target in the reading order
				// and the best path to that target is with the next step
				// being the first in reading order
				tdir = 5
				tloc = xy{10000,10000}
				for i, pth := range pths {
					if pth.vld {
						if pth.loc.y < tloc.y ||
							(pth.loc.y == tloc.y && pth.loc.x < tloc.x) ||
							(pth.loc.y == tloc.y && pth.loc.x < tloc.x && pth.sdir < tdir) {
							tdir = pth.sdir
							tloc = pth.loc
							npix = i
						}
					}
				}

				// move & update the map
				if tdir != 5 {
					_, nloc = neigh(ppl[ppix].loc, tdir, mp)
					mp[pp.loc.y][pp.loc.x] = emp
					delete(pploc, pp.loc)
					ppl[ppix].loc = nloc
					if pp.elf {
						mp[nloc.y][nloc.x] = elf
					} else {
						mp[nloc.y][nloc.x] = gob
					}
					pploc[nloc] = ppl[ppix]
				}
			}

			// attacking is simple since I have all data
			if attack != nil {
				if attack.elf {
					attack.hp -= 3
				} else {
					attack.hp -= strength
				}
				if attack.hp < 1 {
					// remove from map and loc index but remove from ppl at the end
					// in order to not mess up the people loop we are in
					mp[attack.loc.y][attack.loc.x] = emp
					delete(pploc, attack.loc)
					died = true
				}
			}
		}

		// clean up the ppl list after somebody died
		// and rebuild the pploc location list as the indexing
		// of ppl is changing
		// als detect the ending scenario as it occurse after somebody died
		totHp   = 0 // all health points
		elfCnt  = 0 // the amount of elves
		gobCnt := 0 // the amount of goblins
		if died {

			// create new versions of ppl and pploc
			// and transport only the living ones
			nppl   := make([]*people, 0, len(ppl))
			npploc := make(map[xy]*people)
			for _,pp := range ppl {
				if pp.hp > 0 {
					nppl = append(nppl, pp)
					npploc[pp.loc] = pp
					if pp.elf {
						elfCnt++
					} else {
						gobCnt++
					}
					totHp += pp.hp
				}
			}
			ppl = nppl
			pploc = npploc

			// if one groups is gone, give up
			if (elfCnt == 0) || (gobCnt == 0) {
				break
			}
		}
	}
	return
}

// MAIN ----
func main () {

	// the direction vectors in reading order
	vdir = []xy{xy{0,-1}, xy{-1,0}, xy{1,0}, xy{0,1}}

	dataset := ""
	if len(os.Args) < 2 || os.Args[1] == "" {
		fmt.Println("No argument given - trying 'test' dataset.")
		dataset = "test"
	} else {
		dataset = os.Args[1]
	}

	start   := time.Now()
	input   := readTxtFile("d15." + dataset + ".txt")

	// part 1
	mp, ppl, pploc := parseFile(input)
	_, totHp, cnt  := battle(3, mp, ppl, pploc)
	fmt.Println("\nPart 1\n---------------")
	fmt.Println("Last completed round with attack strength 3: ", cnt)
	fmt.Println("Remaining health of the Goblins: ", totHp)
	fmt.Println("Check sum: ", totHp*(cnt),"\n")

	// part 2
	fmt.Println("Part 2\n---------------")
	mp, ppl, pploc  = parseFile(input) // refresh the input after part 1
	oElfCnt := countElfs(ppl)          // the original elf count
	for i := 4; true; i++ {
		elfCnt, totHp, cnt  := battle(i, mp, ppl, pploc)
		fmt.Printf("Attack strength %2v, Elfs left: (%2v of %2v), Checksum: %6v (HP: %4v * Rnds: %3v)\n", i, elfCnt, oElfCnt, totHp*cnt, totHp, cnt)
		if elfCnt == oElfCnt {
			fmt.Println()
			break
		}
		mp, ppl, pploc = parseFile(input)
	}

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}