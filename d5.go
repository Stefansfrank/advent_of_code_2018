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
	inp   := readTxtFile("d5." + dataset + ".txt")[0]
	oinp  := inp

	cap   := byte('a' - 'A')
	dl    := byte(0)
	found := true

	tmpC  := byte(0)
	lets  := map[byte]bool{}
	
	for found {
		found = false
		for i := 0; i < len(inp)-1; i++ {
			if inp[i+1] > inp[i] {
				dl   = inp[i+1] - inp[i]
				tmpC = inp[i]
			} else {
				dl   = inp[i] - inp[i+1]
				tmpC = inp[i+1]
			}
			if dl == cap {
				inp        = inp[:i] + inp[i+2:]
				lets[tmpC] = true
				found      = true
				break
			}
		}
	}

	fmt.Println(len(inp), "units remain.")

	for c, _ := range lets {

		binp := make([]byte, 0, len(oinp))
		for _, cc := range oinp {
			if c == byte(cc) || c + cap == byte(cc) {
				continue
			} else {
				binp = append(binp, byte(cc))
			}
		}

		found = true
		for found {
			found = false
			for i := 0; i < len(binp)-1; i++ {
				if binp[i+1] > binp[i] {
					dl   = binp[i+1] - binp[i]
					tmpC = binp[i+1]
				} else {
					dl   = binp[i] - binp[i+1]
					tmpC = binp[i+1]
				}
				if dl == cap {
					binp  = append(binp[:i], binp[i+2:]...)
					found = true
					break
				}
			}
		}

		fmt.Printf("Removing all %c elements results in a polymer of length %v.\n", c, len(binp))
	}

 	fmt.Printf("Execution time: %v\n", time.Since(start))
}