
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func scanLayers(x, y int, data io.Reader) [][]int {
	size := x * y

	b := bufio.NewReader(data)
	out := [][]int{}
	index := 0
	this := make([]int, 25 * 6)
	var err error

	for {
		c, err := b.ReadByte()
		if err != nil {
			break
		}
		if index == size {
			th_copy := make([]int, len(this))
			copy(th_copy, this)
			out = append(out, th_copy)
			this = make([]int, 25 * 6)
			index = 0
		}
		num, err := strconv.Atoi(string(c))
		if err != nil {
			continue
		}
		this[index] = num
		index++
	}

	if err != nil && err != io.EOF {
		panic(err)
	}

	return out
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func main() {
	file, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	layers := scanLayers(25, 6, file)

	topLayer := make([]int, 25 * 6)
	for i := range topLayer {
		topLayer[i] = 2;
	}

	for _, layer := range layers {
		for i, p := range layer {
			if topLayer[i] != 2 {
				continue
			}
			if p == 2 {
				continue
			}
			topLayer[i] = p
		}
	}

	for y := 0; y < 6; y++ {
		for x := 0; x < 25; x++ {
			if topLayer[y * 25 + x] == 0 {
				fmt.Print(" ")
			} else {
				fmt.Print("*")
			}
			// fmt.Print(topLayer[y * 25 + x])
		}
		fmt.Println()
	}
}
