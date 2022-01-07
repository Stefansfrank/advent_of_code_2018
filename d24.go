package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
	"regexp"
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

const imm = 0 // immune system
const inf = 1 // infection
var   typ map[string]byte // mapping the strings of the attack types to the associated bit
var  rtyp map[byte]string // mapping the bit to the string of the attack types

// the representation of a group
type grp struct {
	trp int   // the troup it belongs to (immune system or infection)
	num int   // number of units
	hp  int   // hit points per unit
	str int   // attack strength per unit
	typ byte  // attack type
	inv int   // initiative
	wk  byte  // weaknesses (bit coded)
	imm byte  // immunity (bit coded)
	// --- vvv filled during target selection phase vvv ---
	att *grp  // the group to be attacked
	def bool  // am I already targeted?
}

// input parser using Regex computes both the army and initiative sorted index to the army as the latter never changes
func parseFile (lines []string, boost int) (army [2][]*grp, invIx []*grp) {
	typ  = map[string]byte{ "bludgeoning":byte(1), "fire":byte(2), "cold":byte(4), "slashing":byte(8), "radiation":byte(16) }
	rtyp = map[byte]string{ byte(1):"bludgeoning", byte(2):"fire", byte(4):"cold", byte(8):"slashing", byte(16):"radiation" }
	army[0] = []*grp{}
	army[1] = []*grp{}
	invIx   = []*grp{}
	var curTrp int

	reUn  := regexp.MustCompile(`([0-9]+) units each with ([0-9]+) hit [a-z ,;()]+ does ([0-9]+) ([a-z]+) damage at initiative ([0-9]+)`)
	reIm  := regexp.MustCompile(`[(;] ?immune to ([a-z]+)(?:, )?([a-z]+)?(?:, )?([a-z]+)?(?:, )?([a-z]+)?(?:, )?([a-z]+)?[;)]`)
	reWk  := regexp.MustCompile(`[(;] ?weak to ([a-z]+)(?:, )?([a-z]+)?(?:, )?([a-z]+)?(?:, )?([a-z]+)?(?:, )?([a-z]+)?[;)]`)

	for _, line := range lines {
		if line == "Immune System:" {
			curTrp = imm
		} else if line == "Infection:" {
			curTrp = inf
		} else if len(line) > 0 {

			match := reUn.FindStringSubmatch(line)
			g     := grp{ trp: curTrp, num: atoi(match[1]), hp: atoi(match[2]), str: atoi(match[3]), typ: typ[match[4]], inv: atoi(match[5]) }
			if curTrp == imm {
				g.str += boost
			}

			match  = reIm.FindStringSubmatch(line)
			for _, m := range match {
				g.imm |= typ[m]
			}

			match  = reWk.FindStringSubmatch(line)
			for _, m := range match {
				g.wk |= typ[m]
			}

			army[curTrp] = append(army[curTrp], &g)
			invIx = append(invIx, &g)
		}
	}

	sort.Slice(invIx, func (i,j int) bool { return invIx[i].inv > invIx[j].inv })
	return
}

// effective power of a group
func (u grp) eff() int {
	return u.num * u.str
}

// executes one round of fighting. Besides the army itself, it takes the index table listing the groups sorted by initiative
func oneFight(army [2][]*grp, invIx []*grp) {

	// build effective power sorted index (preserving the initiative sorting for equal effective power)
	effIx := make([]*grp, len(invIx))
	copy(effIx, invIx)
	sort.SliceStable(effIx, func (i, j int) bool { return effIx[i].eff() > effIx[j].eff() })

	// target selection
	for _, grp := range effIx {

		var dm int
		dmMax := 0
		for _, tgrp := range effIx {

			// exclude target if it either belongs to the same army or is already defending or is immune to attack type or is already dead
			if grp.trp == tgrp.trp || tgrp.def || tgrp.imm & grp.typ > 0 || tgrp.num < 1 {
				continue
			}

			// the way the sorting is organized, we will look at troups with higher effective power first
			// and in equality cases we'll look at higher initiative troups first so if we take the first
			// troup with highest damage, we properly deal with equality cases
			dm = grp.eff() * (1 + int(tgrp.wk & grp.typ/grp.typ))
			if dm > dmMax {
				grp.att  = tgrp
				dmMax  = dm
			}
		}

		if grp.att != nil {
			grp.att.def = true
		}
	}

	// fight (since I represented immunities and weaknesses as bits, a simple & operation allows me to avoid a condition)
	for _, grp := range invIx {
		if grp.att != nil {
			grp.att.num -= grp.eff() * (1 + int(grp.att.wk & grp.typ/grp.typ)) / grp.att.hp
			grp.att.num = max(grp.att.num, 0)
			grp.att.def = false
			grp.att = nil 
		}
	}

	return
}

// debugging - print out the current state of the armies
func dump(army [2][]*grp) {
	fmt.Println("\nImmune System\n-------------")
	for _,g := range army[0] {
		g.dump(true)
	}
	fmt.Println("\nInfection    \n-------------")
	for _,g := range army[1] {
		g.dump(true)
	}
	fmt.Println()
}

// debugging - print out one group ('exp' controls whether the to-be-attacked group is detailed out)
func (g *grp) dump(exp bool) {
	fmt.Printf("%v units (hp: %v, str: %v, typ: %v, inv: %v, imm: %v, wk: %v) - defending: %v",
					g.num, g.hp, g.str, rtyp[g.typ], g.inv, g.imm, g.wk, g.def)
	if g.att != nil {
		if exp {
			fmt.Print("\nTarget: ")
			g.att.dump(false)
		} else {
			fmt.Println(" Target: true")
		}
	} else {
		fmt.Println(" Target: false")
	}
}

// counts units of both armies and indicates who is winning
// returns minimum count, maximum count, win flag (from the reindeer's perspective)
func count(army [2][]*grp) (int, int, bool) {
	var sm0,sm1 int
	for _, g := range army[0] {
		sm0 += g.num
	}
	for _, g := range army[1] {
		sm1 += g.num
	}
	if sm0 < sm1 {
		return sm0, sm1, false
	}
	return sm1, sm0, true
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
	input := readFileStr("d24." + dataset + ".txt")

	// part 1 - repeat fights until the minimum count of armies is 0
	army, invIx := parseFile(input, 0)
	minCnt := 1
	maxCnt := 0
	var win bool
	for minCnt > 0 {
		minCnt, maxCnt, win = count(army)
		oneFight(army, invIx)
	}

	fmt.Printf("\n%v groups left at the end of the fight - ", maxCnt)
	if win {
		fmt.Println("Reindeer lives !!\n")
	} else {
		fmt.Println("Reindeer dead :'(\n")
	}

	// part 2 - loop through boost values until a win comes up
	// note that stalemates can occur if the only armies left 
	// have immunity to all attacks the other army can still do
	for i := 1; i < 1000; i++ {
		army, invIx = parseFile(input, i)
		minCnt = 1
		pmin  := 1
		pmax  := 0
		stale := false
		for minCnt > 0 {
			minCnt, maxCnt, win = count(army)
			if pmin == minCnt && pmax == maxCnt {
				stale = true
				break
			} else {
				pmin = minCnt
				pmax = maxCnt
			}
			oneFight(army, invIx)
		}
		if stale {
			// fmt.Printf("With boost %v, the battle ended in a stalemate\n", i)
		} else if win {
			fmt.Printf("With boost %v, the reindeer lives and has %v immune units left!\n", i, maxCnt)
			break
		} else {
			// fmt.Printf("With boost %v, the reindeer is dead :( - the infection still has %v units.\n", i, maxCnt)
		}
	}
	
 	fmt.Printf("\nExecution time: %v\n", time.Since(start))
}