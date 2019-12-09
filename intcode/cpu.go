package main

import (
	"fmt"
	"math/big"
	"sync"
)

func ics(v string) big.Int {
	var x big.Int
	err := x.UnmarshalText([]byte(v))
	if err != nil {
		panic(err)
	}
	return x
}

func ici(i int) big.Int {
	return *big.NewInt(int64(i))
}

func toi(i *big.Int) int {
	if i.IsInt64() {
		return int(i.Int64())
	} else {
		panic("impossible conversion to int")
	}
}

func tos(i *big.Int) string {
	b, err := i.MarshalText()
	if err != nil {
		panic(err)
	}
	return string(b)
}

type IntcodeCPU struct {
	Memory         []big.Int
	RelativeBase   int
	ProgramCounter int
}

func ExecuteProgram(cpu IntcodeCPU, input <-chan big.Int, output chan<- big.Int) {
interp:
	for {
		did_jump := false
		in, consumed := ScanInstruction(cpu.Memory[cpu.ProgramCounter:])
		fmt.Println(in)
		switch in.Opcode {
		case OP_ADD, OP_MULTIPLY:
			v1 := in.Args[0].Value(cpu)
			v2 := in.Args[1].Value(cpu)
			var out big.Int

			if in.Opcode == OP_ADD {
				out.Add(&v1, &v2) // Why that's out = v1 + v2 is beyond me
			} else if in.Opcode == OP_MULTIPLY {
				out.Mul(&v1, &v2)
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
			itarget := toi(&target)
			if (in.Opcode == OP_JUMP_T && toi(&v) != 0) || (in.Opcode == OP_JUMP_F && toi(&v) == 0) {
				cpu.ProgramCounter = itarget
				did_jump = true
			}
		case OP_LESS, OP_EQUALS:
			c1 := in.Args[0].Value(cpu)
			c2 := in.Args[1].Value(cpu)
			p := in.Args[2].PValue(cpu)
			if (in.Opcode == OP_LESS && c1.Cmp(&c2) == -1) || (in.Opcode == OP_EQUALS && c1.Cmp(&c2) == 0) {
				*p = ici(1)
			} else {
				*p = ici(0)
			}
		case OP_SETBASE:
			base := in.Args[0].Value(cpu)
			ibase := toi(&base)
			cpu.RelativeBase += ibase
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

func main() {
	program := ReadProgram("input")

	// fmt.Println(program)

	/*
		instrs := ScanInstructions(program)
		for _, v := range instrs {
			fmt.Println(v)
		}
	*/
	input := make(chan big.Int)
	output := make(chan big.Int)

	memory := make([]big.Int, 100000)
	copy(memory, program)

	cpu := IntcodeCPU{
		Memory:       memory,
		RelativeBase: 0,
		ProgramCounter: 0,
	}

	go ExecuteProgram(cpu, input, output)

	var inp big.Int
	fmt.Print("input: ")
	fmt.Scanf("%d", &inp)
	input <- inp
	close(input)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for x := range output {
			fmt.Println("output:", tos(&x))
		}
		wg.Done()
	}()
	wg.Wait()
}
