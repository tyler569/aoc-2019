package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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

type Program []int

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

func readProgram(filename string) Program {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	b := bufio.NewScanner(f)
	b.Split(scanCommas)

	program := make(Program, 0, 100)

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

	PSOP_DATA   Op = 0
	OP_ADD      Op = 1
	OP_MULTIPLY Op = 2
	OP_INPUT    Op = 3
	OP_OUTPUT   Op = 4
	OP_JUMP_T   Op = 5
	OP_JUMP_F   Op = 6
	OP_LESS     Op = 7
	OP_EQUALS   Op = 8
	OP_HALT     Op = 99
)

var modeStrings map[Mode]string = map[Mode]string{
	MODE_POSITION:  "Position",
	MODE_IMMEDIATE: "Immediate",
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
}

func (o Op) String() string {
	return opStrings[o]
}

type Arg struct {
	v    int
	mode Mode
}

func (a Arg) String() string {
	return fmt.Sprintf("(%v %v)", a.mode, a.v)
}

type Instruction struct {
	opcode Op
	args   []Arg
}

func (i Instruction) String() string {
	if len(i.args) > 0 {
		return fmt.Sprintf("%v %v", i.opcode, i.args)
	} else {
		return fmt.Sprintf("%v", i.opcode)
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
}

func scanInstruction(program []int) (i Instruction, consume int) {
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

func scanInstructions(program []int) []Instruction {
	var code []Instruction
	var index int

	for index < len(program) {
		i, consume := scanInstruction(program[index:])
		code = append(code, i)
		index += consume
	}

	return code
}

func (a Arg) Value(p Program) int {
	if a.mode == MODE_IMMEDIATE {
		return a.v
	} else if a.mode == MODE_POSITION {
		return p[a.v]
	}
	panic("invalid argument mode")
}

func (a Arg) PValue(p Program) *int {
	if a.mode == MODE_POSITION {
		return &p[a.v]
	}
	panic("invalid pointer value")
}

func executeProgram(program Program, input <-chan int, output chan<- int) {
	// code := scanInstructions(p)

	var index int
	interp: for {
		did_jump := false
		in, consumed := scanInstruction(program[index:])
		// fmt.Println(in)
		switch in.opcode {
		case OP_ADD, OP_MULTIPLY:
			v1 := in.args[0].Value(program)
			v2 := in.args[1].Value(program)
			var out int

			if in.opcode == OP_ADD {
				out = v1 + v2
			} else if in.opcode == OP_MULTIPLY {
				out = v1 * v2
			}

			*(in.args[2].PValue(program)) = out
		case OP_OUTPUT:
			v := in.args[0].Value(program)
			output <- v
		case OP_INPUT:
			v := <-input
			ptr := in.args[0].PValue(program)
			*ptr = v
		case OP_HALT:
			close(output)
			break interp
		case OP_JUMP_T, OP_JUMP_F:
			v := in.args[0].Value(program)
			target := in.args[1].Value(program)
			if (in.opcode == OP_JUMP_T && v != 0) || (in.opcode == OP_JUMP_F && v == 0) {
				index = target
				did_jump = true
			}
		case OP_LESS, OP_EQUALS:
			c1 := in.args[0].Value(program)
			c2 := in.args[1].Value(program)
			p := in.args[2].PValue(program)
			if (in.opcode == OP_LESS && c1 < c2) || (in.opcode == OP_EQUALS && c1 == c2) {
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

func permutations() [][]int {
	r := [][]int{}
	bot := 5
	top := 10
	for i := bot; i < top; i++ {
		for j := bot; j < top; j++ {
			if j == i { continue }
			for k := bot; k < top; k++ {
				if k == j || k == i { continue }
				for l := bot; l < top; l++ {
					if l == k || l == j || l == i { continue }
					for m := bot; m < top; m++ {
						if m == l || m == k || m == j || m == i { continue }
						r = append(r, []int{i, j, k, l, m})
					}
				}
			}
		}
	}

	return r
}

func createAmplifier(program []int, input <-chan int, phase int) <-chan int {
	output := make(chan int)
	localInput := make(chan int, 2)
	localProgram := make(Program, len(program))
	copy(localProgram, program)

	localInput <- phase

	go func() {
		for x := range input {
			localInput <- x
		}
		close(localInput)
	}()

	go executeProgram(localProgram, localInput, output)
	return output
}

func main() {
	program := readProgram("input")

	// fmt.Println(program)

	/*
	instrs := scanInstructions(program)
	for _, v := range instrs {
		fmt.Println(v)
	}
	*/

	perms := permutations()
	max := -1

	for _, permutation := range perms {
		fmt.Println("trying", permutation)

		outerInput := make(chan int, 2)
		outerInput <- 0

		var lastOutput <-chan int = outerInput

		for _, x := range permutation {
			lastOutput = createAmplifier(program, lastOutput, x)
		}

		var finalResult int
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			for x := range lastOutput {
				finalResult = x
				fmt.Println("looping", x)
				outerInput <- x
			}
			wg.Done()
		}()

		wg.Wait()

		fmt.Println("yields:", finalResult)

		if finalResult > max {
			max = finalResult
		}
	}

	fmt.Println(max)
}

