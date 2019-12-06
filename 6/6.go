
package main

import (
	"fmt"
	"os"
)

func splitParen(s string) (string, string) {
	for i := 0; i < len(s); i++ {
		if s[i] == ')' {
			return s[:i], s[i+1:]
		}
	}
	panic("invalid string in input file")
}

type OrbitMap map[string]string

func countToCOM(s string, m OrbitMap) int {
	if s == "COM" {
		return 0
	} else {
		return countToCOM(m[s], m) + 1
	}
}

func countOrbits(m OrbitMap) int {
	orbits := 0

	for child := range m {
		orbits += countToCOM(child, m)
	}
	return orbits
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	orbitMap := make(OrbitMap)

	for {
		var str string
		_, err := fmt.Fscanln(f, &str)
		if err != nil {
			fmt.Println(err)
			break
		}

		parent, child := splitParen(str)
		orbitMap[child] = parent
	}

	// fmt.Println(orbitMap)

	fmt.Println(countOrbits(orbitMap))
}

