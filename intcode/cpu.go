package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

type IntcodeCPU struct {
	Memory         []int
	RelativeBase   int
	ProgramCounter int
}

func ExecuteProgram(cpu IntcodeCPU, input <-chan int, output chan<- int) {
interp:
	for {
		did_jump := false
		in, consumed := ScanInstruction(cpu.Memory[cpu.ProgramCounter:])
		// fmt.Println(in)
		switch in.Opcode {
		case OP_ADD, OP_MULTIPLY:
			v1 := in.Args[0].Value(cpu)
			v2 := in.Args[1].Value(cpu)
			var out int

			if in.Opcode == OP_ADD {
				out = v1 + v2
			} else if in.Opcode == OP_MULTIPLY {
				out = v1 * v2
			}

			*(in.Args[2].PValue(cpu)) = out
		case OP_OUTPUT:
			v := in.Args[0].Value(cpu)
			output <- v
		case OP_INPUT:
			v := <-input
			ptr := in.Args[0].PValue(cpu)
			*ptr = v
		case OP_HALT:
			close(output)
			break interp
		case OP_JUMP_T, OP_JUMP_F:
			v := in.Args[0].Value(cpu)
			target := in.Args[1].Value(cpu)
			if (in.Opcode == OP_JUMP_T && v != 0) || (in.Opcode == OP_JUMP_F && v == 0) {
				cpu.ProgramCounter = target
				did_jump = true
			}
		case OP_LESS, OP_EQUALS:
			c1 := in.Args[0].Value(cpu)
			c2 := in.Args[1].Value(cpu)
			p := in.Args[2].PValue(cpu)
			if (in.Opcode == OP_LESS && c1 < c2) || (in.Opcode == OP_EQUALS && c1 == c2) {
				*p = 1
			} else {
				*p = 0
			}
		case OP_SETBASE:
			base := in.Args[0].Value(cpu)
			cpu.RelativeBase += base
		default:
			fmt.Printf("Invalid opcode: %d at index %d\n",
				cpu.Memory[cpu.ProgramCounter], cpu.ProgramCounter)
			panic("cannot continue")
		}

		if !did_jump {
			cpu.ProgramCounter += consumed
		}
	}
}

func displayHelp() {
	fmt.Println("intcode: an intcode CPU")
	fmt.Println("  INPUT instructions are read from STDIN")
	fmt.Println("  OUTPUT instructions are written to STDOUT")
	fmt.Println(" --program FILE       The file with the program")
}

func main() {
	var file = flag.String("program", "input", "program file")
	flag.Parse()

	fmt.Fprintln(os.Stderr, "Running intcode")
	program := ReadProgram(*file)

	// fmt.Println(program)

	/*
		instrs := ScanInstructions(program)
		for _, v := range instrs {
			fmt.Println(v)
		}
	*/
	input := make(chan int)
	output := make(chan int)

	memory := make([]int, 100000)
	copy(memory, program)

	cpu := IntcodeCPU{
		Memory:       memory,
		RelativeBase: 0,
		ProgramCounter: 0,
	}

	go ExecuteProgram(cpu, input, output)

	go func() {
		for {
			var inp int
			fmt.Scanf("%d", &inp)
			input <- inp
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for x := range output {
			fmt.Println(x)
		}
		wg.Done()
	}()
	wg.Wait()
}
