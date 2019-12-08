
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
	this := []int{}
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
			this = []int{}
			index = 0
		}
		num, err := strconv.Atoi(string(c))
		if err != nil {
			continue
		}
		this = append(this, num)
		index++
	}

	if err != nil && err != io.EOF {
		panic(err)
	}

	return out
}

func countNs(layer []int, v int) int {
	var count int
	for _, l := range layer {
		if l == v {
			count++
		}
	}
	return count
}

func main() {
	file, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	layers := scanLayers(25, 6, file)

	minZeros := 999
	var min []int
	for _, l := range layers {
		zeros := countNs(l, 0)
		if zeros < minZeros {
			minZeros = zeros
			min = l
		}
	}

	fmt.Println(min)

	ones := countNs(min, 1)
	twos := countNs(min, 2)

	fmt.Println(ones, twos, ones * twos)
}
