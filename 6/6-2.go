
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

func pathToCOM(s string, m OrbitMap) []string {
	if s == "COM" {
		return []string{}
	} else {
		return append(pathToCOM(m[s], m), s)
	}
}

func lengthOfCommonPrefix(a1, a2 []string) int {
	for i := 0; ; i++ {
		if i >= len(a1) || i >= len(a2) {
			return i
		}
		if a1[i] == a2[i] {
			continue
		}
		return i
	}
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

	youPath := pathToCOM("YOU", orbitMap)
	sanPath := pathToCOM("SAN", orbitMap)

	/*
	fmt.Println(youPath)
	fmt.Println(sanPath)
	*/

	common := lengthOfCommonPrefix(youPath, sanPath)
	transfers := len(youPath) + len(sanPath) - 2*common - 2

	fmt.Println(transfers)
}

