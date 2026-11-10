package main

import (
	"errors"
	"strconv"
	"strings"
)

// Farm holds the parsed ant farm.
type Farm struct {
	ants  int
	names []string       // room names by index
	index map[string]int // room name -> index
	links [][]int        // neighbours of each room
	start int
	end   int
}

func parse(text string) (*Farm, error) {
	f := &Farm{index: map[string]int{}, start: -1, end: -1}
	lines := strings.Split(strings.ReplaceAll(text, "\r", ""), "\n")
	seenLink := map[string]bool{}
	gotAnts := false
	next := "" // "start" or "end" if the next room is special

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "##start" || line == "##end" {
			if next != "" {
				return nil, errors.New("bad ##start/##end command")
			}
			next = line[2:]
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue // comment or unknown command
		}
		if !gotAnts {
			n, err := strconv.Atoi(line)
			if err != nil || n <= 0 {
				return nil, errors.New("invalid number of ants")
			}
			f.ants = n
			gotAnts = true
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 3 {
			name := parts[0]
			if strings.HasPrefix(name, "L") {
				return nil, errors.New("invalid room name: " + name)
			}
			if _, err := strconv.Atoi(parts[1]); err != nil {
				return nil, errors.New("invalid coordinates for room " + name)
			}
			if _, err := strconv.Atoi(parts[2]); err != nil {
				return nil, errors.New("invalid coordinates for room " + name)
			}
			if strings.Contains(name, "-") {
				return nil, errors.New("invalid room name: " + name)
			}
			if _, ok := f.index[name]; ok {
				return nil, errors.New("duplicated room: " + name)
			}
			f.index[name] = len(f.names)
			f.names = append(f.names, name)
			f.links = append(f.links, nil)
			if next == "start" {
				if f.start != -1 {
					return nil, errors.New("more than one start room")
				}
				f.start = f.index[name]
			} else if next == "end" {
				if f.end != -1 {
					return nil, errors.New("more than one end room")
				}
				f.end = f.index[name]
			}
			next = ""
			continue
		}
		if next != "" {
			return nil, errors.New("##" + next + " is not followed by a room")
		}
		names := strings.Split(line, "-")
		if len(parts) != 1 || len(names) != 2 {
			return nil, errors.New("invalid line: " + line)
		}
		a, okA := f.index[names[0]]
		b, okB := f.index[names[1]]
		if !okA || !okB {
			return nil, errors.New("link to unknown room: " + line)
		}
		if a == b {
			return nil, errors.New("room links to itself: " + line)
		}
		key := names[0] + "-" + names[1]
		if a > b {
			key = names[1] + "-" + names[0]
		}
		if seenLink[key] {
			return nil, errors.New("duplicated link: " + line)
		}
		seenLink[key] = true
		f.links[a] = append(f.links[a], b)
		f.links[b] = append(f.links[b], a)
	}

	if !gotAnts {
		return nil, errors.New("invalid number of ants")
	}
	if f.start == -1 {
		return nil, errors.New("no start room found")
	}
	if f.end == -1 {
		return nil, errors.New("no end room found")
	}
	return f, nil
}
