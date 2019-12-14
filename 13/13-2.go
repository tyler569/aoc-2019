package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"flag"
	"sync"
)

func min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func scanCommas(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data); i += 1 {
		if data[i] == ',' {
			return i + 1, data[:i], nil
		}
	}
	if !atEOF {
		return 0, nil, nil
	}
	return 0, data, bufio.ErrFinalToken
}

func ReadProgram(filename string) []int {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	b := bufio.NewScanner(f)
	b.Split(scanCommas)

	program := make([]int, 0, 100)

	for b.Scan() {
		num := strings.TrimSpace(b.Text())
		val, err := strconv.Atoi(num)
		if err != nil {
			panic(err)
		}
		program = append(program, val)
	}
	return program
}

type Op int
type Mode int

const (
	MODE_POSITION  Mode = 0
	MODE_IMMEDIATE Mode = 1
	MODE_RELATIVE  Mode = 2

	PSOP_DATA   Op = 0
	OP_ADD      Op = 1
	OP_MULTIPLY Op = 2
	OP_INPUT    Op = 3
	OP_OUTPUT   Op = 4
	OP_JUMP_T   Op = 5
	OP_JUMP_F   Op = 6
	OP_LESS     Op = 7
	OP_EQUALS   Op = 8
	OP_SETBASE  Op = 9
	OP_HALT     Op = 99
)

var modeStrings map[Mode]string = map[Mode]string{
	MODE_POSITION:  "Position",
	MODE_IMMEDIATE: "Immediate",
	MODE_RELATIVE:  "Relative",
}

func (m Mode) String() string {
	return modeStrings[m]
}

var opStrings map[Op]string = map[Op]string{
	PSOP_DATA:   "data",
	OP_ADD:      "add",
	OP_MULTIPLY: "mul",
	OP_INPUT:    "input",
	OP_OUTPUT:   "output",
	OP_HALT:     "halt",
	OP_JUMP_T:   "jumpt",
	OP_JUMP_F:   "jumpf",
	OP_LESS:     "less",
	OP_EQUALS:   "greater",
	OP_SETBASE:  "setbase",
}

func (o Op) String() string {
	return opStrings[o]
}

type Arg struct {
	V    int
	Mode Mode
}

func (a Arg) String() string {
	return fmt.Sprintf("(%v %v)", a.Mode, a.V)
}

type Instruction struct {
	Opcode Op
	Args   []Arg
}

func (i Instruction) String() string {
	if len(i.Args) > 0 {
		return fmt.Sprintf("%v %v", i.Opcode, i.Args)
	} else {
		return fmt.Sprintf("%v", i.Opcode)
	}
}

func opDecode(op int) (Op, [3]Mode) {
	return Op(op % 100),
		[3]Mode{
			Mode((op / 100) % 10),
			Mode((op / 1000) % 10),
			Mode((op / 10000) % 10),
		}
}

var opArgCount map[Op]int = map[Op]int{
	OP_ADD:      3,
	OP_MULTIPLY: 3,
	OP_INPUT:    1,
	OP_OUTPUT:   1,
	OP_HALT:     0,
	OP_JUMP_T:   2,
	OP_JUMP_F:   2,
	OP_LESS:     3,
	OP_EQUALS:   3,
	OP_SETBASE:  1,
}

func ScanInstruction(program []int) (i Instruction, consume int) {
	op, m := opDecode(program[0])

	arg_count, ok := opArgCount[op]
	if !ok {
		return Instruction{PSOP_DATA, []Arg{Arg{program[0], MODE_IMMEDIATE}}}, 1
	}

	var args []Arg

	for i := 0; i < arg_count; i++ {
		args = append(args, Arg{program[i+1], m[i]})
	}

	return Instruction{op, args}, arg_count + 1
}

func ScanInstructions(program []int) []Instruction {
	var code []Instruction
	var index int

	for index < len(program) {
		i, consume := ScanInstruction(program[index:])
		code = append(code, i)
		index += consume
	}

	return code
}

func (a Arg) Value(p IntcodeCPU) int {
	if a.Mode == MODE_IMMEDIATE {
		return a.V
	} else if a.Mode == MODE_POSITION {
		return p.Memory[a.V]
	} else if a.Mode == MODE_RELATIVE {
		return p.Memory[p.RelativeBase + a.V]
	}
	panic("invalid argument mode")
}

func (a Arg) PValue(p IntcodeCPU) *int {
	if a.Mode == MODE_POSITION {
		return &p.Memory[a.V]
	} else if a.Mode == MODE_RELATIVE {
		return &p.Memory[p.RelativeBase + a.V]
	}
	panic("invalid pointer value")
}

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
			<-input
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
	program := ReadProgram("input2")

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
			input <- 0 // blocker

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
