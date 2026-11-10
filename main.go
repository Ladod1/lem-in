package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: usage: go run . <file>")
		return
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("ERROR: could not read file")
		return
	}
	farm, err := parse(string(data))
	if err != nil {
		fmt.Println("ERROR: invalid data format, " + err.Error())
		return
	}
	paths := findPaths(farm)
	if len(paths) == 0 {
		fmt.Println("ERROR: invalid data format, no path between start and end")
		return
	}
	fmt.Println(strings.TrimRight(string(data), "\n"))
	fmt.Println()
	for _, line := range moveAnts(farm, paths) {
		fmt.Println(line)
	}
}
