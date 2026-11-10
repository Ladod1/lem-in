package main

import (
	"fmt"
	"strings"
)

// moveAnts returns one line per turn with the moves of the ants.
func moveAnts(f *Farm, paths [][]int) []string {
	chosen := assign(paths, f.ants)
	order := make([]int, f.ants) // how many ants went before on the same path
	count := make([]int, len(paths))
	for a, p := range chosen {
		order[a] = count[p]
		count[p]++
	}
	var lines []string
	for t := 0; ; t++ {
		var moves []string
		for a, p := range chosen {
			pos := t - order[a]
			if pos >= 0 && pos < len(paths[p]) {
				moves = append(moves, fmt.Sprintf("L%d-%s", a+1, f.names[paths[p][pos]]))
			}
		}
		if len(moves) == 0 {
			return lines
		}
		lines = append(lines, strings.Join(moves, " "))
	}
}
