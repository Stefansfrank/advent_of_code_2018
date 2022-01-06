package main

import (
	"fmt"
	"strconv"
	"os"
	"bufio"
	"time"
)

// no error handling ...
func readTxtFileInt (name string) (lines []int) {	
	file, _ := os.Open(name)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {		
		lines = append(lines, atoi(scanner.Text()))
	}
	return
}

// inline capable / no error handling Atoi
func atoi (x string) int {
	y, _ := strconv.Atoi(x)
	return y
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
	input  := readTxtFileInt("d1." + dataset + ".txt")

	// Part 1
	res    := 0
	for _, i := range input {
		res += i
	}
	fmt.Println("Final frequency:", res)

	// Part 2
	res  = 0
	hit := map[int]bool{0:true}
	outer: for {
		for _, i := range input {
			res += i
			if hit[res] {
				fmt.Println("First repeated frequency:", res)
				break outer
			} else {
				hit[res] = true
			}
		}	
	}

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}