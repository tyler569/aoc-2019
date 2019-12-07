package main

import (
	"fmt"
	"sync"

	"github.com/tyler569/aoc-2019/intcode/bits"
)

func executeProgram(program bits.Program, input <-chan int, output chan<- int) {
	var index int
	interp: for {
		did_jump := false
		in, consumed := bits.ScanInstruction(program[index:])
		switch in.Opcode {
		case bits.OP_ADD, bits.OP_MULTIPLY:
			v1 := in.Args[0].Value(program)
			v2 := in.Args[1].Value(program)
			var out int

			if in.Opcode == bits.OP_ADD {
				out = v1 + v2
			} else if in.Opcode == bits.OP_MULTIPLY {
				out = v1 * v2
			}

			*(in.Args[2].PValue(program)) = out
		case bits.OP_OUTPUT:
			v := in.Args[0].Value(program)
			output <- v
		case bits.OP_INPUT:
			v := <-input
			ptr := in.Args[0].PValue(program)
			*ptr = v
		case bits.OP_HALT:
			close(output)
			break interp
		case bits.OP_JUMP_T, bits.OP_JUMP_F:
			v := in.Args[0].Value(program)
			target := in.Args[1].Value(program)
			if (in.Opcode == bits.OP_JUMP_T && v != 0) || (in.Opcode == bits.OP_JUMP_F && v == 0) {
				index = target
				did_jump = true
			}
		case bits.OP_LESS, bits.OP_EQUALS:
			c1 := in.Args[0].Value(program)
			c2 := in.Args[1].Value(program)
			p := in.Args[2].PValue(program)
			if (in.Opcode == bits.OP_LESS && c1 < c2) || (in.Opcode == bits.OP_EQUALS && c1 == c2) {
				*p = 1
			} else {
				*p = 0
			}
		default:
			fmt.Printf("Invalid opcode: %d at index %d\n", program[index], index)
			panic("cannot continue")
		}

		if !did_jump {
			index += consumed
		}
	}
}

func main() {
	program := bits.ReadProgram("input")

	// fmt.Println(program)

	/*
	instrs := bits.ScanInstructions(program)
	for _, v := range instrs {
		fmt.Println(v)
	}
	*/
	input := make(chan int)
	output := make(chan int)

	go executeProgram(program, input, output)

	var inp int
	fmt.Print("input: ")
	fmt.Scanf("%d", &inp)
	input <- inp
	close(input)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for x := range output {
			fmt.Println("output:", x)
		}
		wg.Done()
	}()
	wg.Wait()
}

