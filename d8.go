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

// input parser
func parseFile (lines []string) (nums []int) {

	for _, line := range lines {
		snums := strings.Split(line, " ")
		nums   = make([]int, len(snums))
		for i, snum := range snums {
			nums[i] = atoi(snum)
		}
	}
	return
}

// recursive metadata sum
func metaSum(nums []int, msm int, ix int) (nmsm int, nix int) {
	numNodes := nums[ix] 
	numMeta  := nums[ix+1]

	nix  = ix + 2
	nmsm = msm

	for i := 0; i < numNodes; i++ {
		nmsm, nix = metaSum(nums, nmsm, nix)
	}

	for i := 0; i < numMeta; i++ {
		nmsm += nums[nix]
		nix += 1 	
	}
	return
}

// recursive value search
func value(nums []int, ix int) (val int, nix int) {
	numNodes := nums[ix] 
	numMeta  := nums[ix+1]
	nix  = ix + 2

	if numNodes == 0 {
		for i := 0; i < numMeta; i++ {
			val += nums[nix]
			nix += 1 	
		}	
		return	
	}

	if numMeta > 0 {

		vl := make([]int, numNodes)
		for i := 0; i < numNodes; i++ {
			vl[i], nix = value(nums, nix)
		}

		for i := 0; i < numMeta; i++ {
			if nums[nix] > 0 && nums[nix] <= numNodes {
				val += vl[nums[nix]-1]
			}
			nix += 1 	
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

	start := time.Now()
	input := readTxtFile("d8." + dataset + ".txt")
	nums  := parseFile(input)

	msm,_ := metaSum(nums, 0, 0)
	fmt.Println("Meta sum:", msm)

	val,_ := value(nums, 0)
	fmt.Println("Value of root:", val)

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}