package main

import "testing"

func TestParseErrors(t *testing.T) {
	bad := []string{
		"0\n##start\na 0 0\n##end\nb 1 1\na-b\n",      // no ants
		"2\na 0 0\n##end\nb 1 1\na-b\n",               // no start
		"2\n##start\na 0 0\nb 1 1\na-b\n",             // no end
		"2\n##start\na 0 0\n##end\nb 1 1\na-c\n",      // unknown room
		"2\n##start\na 0 0\n##end\nb 1 1\na-a\n",      // self link
		"2\n##start\na 0 0\n##end\nb x 1\na-b\n",      // bad coords
		"2\n##start\na 0 0\n##end\na 1 1\na-b\n",      // duplicated room
		"2\n##start\na 0 0\n##end\nb 1 1\na-b\nb-a\n", // duplicated link
	}
	for _, s := range bad {
		if _, err := parse(s); err == nil {
			t.Errorf("expected error for:\n%s", s)
		}
	}
}

func TestExample(t *testing.T) {
	f, err := parse("3\n##start\n0 1 0\n##end\n1 5 0\n2 9 0\n3 13 0\n0-2\n2-3\n3-1\n")
	if err != nil {
		t.Fatal(err)
	}
	lines := moveAnts(f, findPaths(f))
	if len(lines) != 5 {
		t.Errorf("expected 5 turns, got %d", len(lines))
	}
}
