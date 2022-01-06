package main

import (
	"fmt"
	"os"
	"bufio"
	"time"
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

// compares two IDs and returns the number of differences
// and the substring of letters that are identical
func compID(id1, id2 string) (com string, del int) {

	for i, c := range id1 {
		if id1[i] == id2[i] {
			com += fmt.Sprintf("%c", c)
		} else {
			del += 1
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

	start  := time.Now()
	input  := readTxtFile("d2." + dataset + ".txt")

	// Part 1
	doubles := 0
	triples := 0
	for _, ln := range input {

		// counts the occurences in a map
		hit := map[rune]int{}
		for _, c := range ln {
			hit[c] += 1
		}

		// checks for the occurence of a double ot triple
		double := false
		triple := false
		for _, v := range hit {
			double = double || v == 2
			triple = triple || v == 3
		}
		if double {
			doubles += 1
		}
		if triple {
			triples += 1
		}

	}
	fmt.Println("Checksum:", doubles * triples)

	// Part 2
	out: for i1, id1 := range input {
		for i2,id2 := range input {

			if i1 == i2 {
				continue
			}

			com, del := compID(id1, id2)
			if del == 1 {
				fmt.Println("ID pair found, common part: ", com)
				break out
			}
		}
	}

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}