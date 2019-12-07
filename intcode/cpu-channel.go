package main

import (
	"fmt"
	"sync"

	"github.com/tyler569/aoc-2019/intcode/bits"
)

func executeProgram(program intcode.Program, input <-chan int, output chan<- int) {
	var index int
	interp: for {
		did_jump := false
		in, consumed := intcode.ScanInstruction(program[index:])
		switch in.Opcode {
		case intcode.OP_ADD, intcode.OP_MULTIPLY:
			v1 := in.Args[0].Value(program)
			v2 := in.Args[1].Value(program)
			var out int

			if in.Opcode == intcode.OP_ADD {
				out = v1 + v2
			} else if in.Opcode == intcode.OP_MULTIPLY {
				out = v1 * v2
			}

			*(in.Args[2].PValue(program)) = out
		case intcode.OP_OUTPUT:
			v := in.Args[0].Value(program)
			output <- v
		case intcode.OP_INPUT:
			v := <-input
			ptr := in.Args[0].PValue(program)
			*ptr = v
		case intcode.OP_HALT:
			close(output)
			break interp
		case intcode.OP_JUMP_T, intcode.OP_JUMP_F:
			v := in.Args[0].Value(program)
			target := in.Args[1].Value(program)
			if (in.Opcode == intcode.OP_JUMP_T && v != 0) || (in.Opcode == intcode.OP_JUMP_F && v == 0) {
				index = target
				did_jump = true
			}
		case intcode.OP_LESS, intcode.OP_EQUALS:
			c1 := in.Args[0].Value(program)
			c2 := in.Args[1].Value(program)
			p := in.Args[2].PValue(program)
			if (in.Opcode == intcode.OP_LESS && c1 < c2) || (in.Opcode == intcode.OP_EQUALS && c1 == c2) {
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
	program := intcode.ReadProgram("input")

	// fmt.Println(program)

	/*
	instrs := intcode.ScanInstructions(program)
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

