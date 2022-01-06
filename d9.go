package main

import (
	"fmt"
	"time"
)

// the marble structure
type mrbl struct {
	no  int   // the value
	lft *mrbl // points to the marble to the left
	rgt *mrbl // points to the marble to the right
}

// runs the game with numPl players until the cutoff marble was played
// returns the scores of all players
func game(numPl, cutoff int) (scrs []int) {
	
	scrs  = make([]int, numPl) 
	curr := &mrbl{no:0}
	curr.lft = curr
	curr.rgt = curr
	var sqlft, sqrgt, new *mrbl
	var pl int

	for i := 1; i <= cutoff; i++ {

		if i % 23 > 0 {
			sqlft = curr.rgt  // determine the marble left of the squeeze
			sqrgt = sqlft.rgt // determine the marble right of the squeeze	

			// squeeze in the new marble
			new = &mrbl{no:i}
			sqlft.rgt = new
			sqrgt.lft = new
			new.lft = sqlft
			new.rgt = sqrgt

			// new current marble
			curr = new

		} else {
			pl = (i-1) % numPl // determine the player ID
			scrs[pl] += i      // add the new marble

			sqrgt = curr.lft.lft.lft.lft.lft.lft // determine the marble right of the marble to be removed
			sqlft = sqrgt.lft.lft 				 // determine the marble left of the marble to be removed
			scrs[pl] += sqlft.rgt.no  			 // add score of removed marble befor it's dropped
			sqlft.rgt = sqrgt					 // cross-link the neigbouring marbles ->
			sqrgt.lft = sqlft 					 // thus the removed marble is no longer in the chain

			// new current marble
			curr = sqrgt
		}

	}

	return
}

// generic max function for []int slices
func max(n []int) (mx int) {
	mx = n[0]
	for i := 1; i < len(n); i++ {
		if n[i] > mx {
			mx = n[i]
		}
	}
	return
}

func main() {
	
	// the first numbers in these slices are my personal input
	// the others are the listed examples
	numPl   := []int{455,9,10,13,17,21,30}
	cutoff  := []int{71223,25,1618,7999,1104,6111,5807}
	dataset := 0

	start := time.Now()

	scrs := game(numPl[dataset], cutoff[dataset])
	fmt.Println("High Score:", max(scrs))

	scrs  = game(numPl[dataset], 100*cutoff[dataset])
	fmt.Println("100x High Score:", max(scrs))

	fmt.Println("Execution time:", time.Since(start))

}