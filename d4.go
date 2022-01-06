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

// observation structure
type obsv struct {
	tim time.Time
	typ int
	grd int
}

// input parser using Regex
func parseFile (lines []string) (obss []obsv) {

	reTime := regexp.MustCompile(`1518-([0-9]+)-([0-9]+) ([0-9]+):([0-9]+)`)
	reGrd  := regexp.MustCompile(`Guard #([0-9]+)`)
	obss    = make([]obsv, len(lines))

	for i, line := range lines {
		match := reTime.FindStringSubmatch(line)
		obss[i] = obsv{}
		obss[i].tim = time.Date(1518, time.Month(atoi(match[1])), atoi(match[2]), 
								atoi(match[3]), atoi(match[4]), 0, 0, time.UTC)
		match = reGrd.FindStringSubmatch(line)
		if match != nil {
			obss[i].typ = -1
			obss[i].grd = atoi(match[1])
		} else if line[19] == 'w' {
			obss[i].typ = 1
		}
	}
	return
}

// counts how often the guard sleeps each minute
func countMins(obss []obsv) map[int][]int {

	grdMin := map[int][]int{}
	grd := 0
	from := 0
	for _, obs := range obss {
		switch obs.typ {
		case -1:
			grd = obs.grd
		case 0:
			from = obs.tim.Minute()
		case 1:
			if grdMin[grd] == nil {
				grdMin[grd] = make([]int, 60)
			}
			for m := from; m < obs.tim.Minute(); m++ {
				grdMin[grd][m] += 1
			}
		}
	}

	return grdMin 
}

// finds the maxima for the two strategies
func findMax(grdMin map[int][]int) (strat1 int, strat2 int) {

	grdSum   := map[int]int{} // the sum of sleep for each guard
	grdMax   := map[int]int{} // the maximum number of sleeps at the same minute per guard
	grdMaxMn := map[int]int{} // the minute the guard sleeps most
	for grd, min := range grdMin {
		max := 0
		for i, m := range min {
			grdSum[grd] += m
			if max < m {
				max = m
				grdMaxMn[grd] = i
				grdMax[grd]   = m
			}
		}
	}

	// determine the guard with the highest amount of sleep
	maxSm  := 0
	grdNo1 := 0
	for grd, sm := range grdSum {
		if sm > maxSm {
			maxSm  = sm
			grdNo1 = grd
		}
	}

	// determine the guard with the highest number of sleeps at the same minute
	tmax   := 0
	grdNo2 := 0
	for grd, max := range grdMax {
		if max > tmax {
			tmax   = max
			grdNo2 = grd
		}
	}

	return grdNo1 * grdMaxMn[grdNo1], grdNo2 * grdMaxMn[grdNo2]
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
	input  := readTxtFile("d4." + dataset + ".txt")
	obss   := parseFile(input)
	sort.Slice(obss, func(i, j int) bool { return obss[i].tim.Before(obss[j].tim) })

	grdMin := countMins(obss)
	s1, s2 := findMax(grdMin)

	fmt.Println("Strategy 1: ", s1)
	fmt.Println("Strategy 2: ", s2)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}