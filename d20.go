package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "time"
    "log"
    "math"
)

// read file into []byte
func readFile (name string) (data []byte) { 
    data, err := ioutil.ReadFile(name)
    if err != nil {
            log.Fatal(err)
    }
    return
}

func min(i, j int) int {
    if i<j {
        return i
    }
    return j
}

func max(i, j int) int {
    if i>j {
        return i
    }
    return j
}

type xy struct {
    x,y int
}

// this assumes that the shortest path to each room is in the input
// the solution would not work for all inputs but the given problems and inputs are 
// constructed in a way that this simple code building the amount of doors from just
// remembering the distance while traversing the input works
func iterate(data []byte) (lng [][]int) {

    vdir  := map[byte]xy{ 'N':xy{0, -1}, 'E':xy{1,0}, 'S':xy{0,1}, 'W':xy{-1,0}}

    // assumes a square-ish shape
    tdim := int(math.Sqrt(float64(len(data)))) + 10 
 
    // the field of distances 
    lng   = make([][]int, tdim)
    for i := 0; i < tdim; i++ {
        lng[i] = make([]int, tdim)
    }

    // the stack remembering the location when brnching
    stack := []xy{}

    // traverse the input
    loc  := xy{tdim/2, tdim/2}
    prv  := loc
    for i := 1; i < len(data) - 1; i++ {
        switch data[i] {
        case '(':
            stack = append(stack, loc)
        case ')':
            loc = stack[len(stack)-1]
            stack = stack[:len(stack)-1]
        case '|':
            loc = stack[len(stack)-1]
        default:
            loc.x += vdir[data[i]].x
            loc.y += vdir[data[i]].y
            if lng[loc.y][loc.x] == 0 {
                lng[loc.y][loc.x] = lng[prv.y][prv.x] + 1
            } else {

                // previously visited
                lng[loc.y][loc.x] = min(lng[loc.y][loc.x], lng[prv.y][prv.x ]+ 1)
            }
        }
        prv = loc
    }

    return
}

// counting the required outputs for both parts
func count(lng [][]int) (mx int, cnt int) {

    for _,rw := range lng {
        for _,n := range rw {
            mx = max(mx, n)
            if n >= 1000 {
                cnt += 1
            }
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
    data  := readFile("d20." + dataset + ".txt")

    lng := iterate(data)
    mx, cnt := count(lng)
    fmt.Printf("\nYou have to pass %v doors for the furthest room\n", mx)
    fmt.Printf("There are %v rooms further than 1000 doors from you\n\n", cnt)

    fmt.Printf("Execution time: %v\n", time.Since(start))
}


