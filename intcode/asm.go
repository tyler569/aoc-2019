package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	// "strconv"
	"strings"
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

func ReadProgram(filename string) []big.Int {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	b := bufio.NewScanner(f)
	b.Split(scanCommas)

	program := make([]big.Int, 0, 100)

	for b.Scan() {
		num := strings.TrimSpace(b.Text())
		val := ics(num)
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
	V    big.Int
	Mode Mode
}

func (a Arg) String() string {
	return fmt.Sprintf("(%v %v)", a.Mode, tos(&a.V))
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

func ScanInstruction(program []big.Int) (i Instruction, consume int) {
	value := toi(&program[0])
	op, m := opDecode(value)

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

func ScanInstructions(program []big.Int) []Instruction {
	var code []Instruction
	var index int

	for index < len(program) {
		i, consume := ScanInstruction(program[index:])
		code = append(code, i)
		index += consume
	}

	return code
}

func (a Arg) Value(p IntcodeCPU) big.Int {
	if a.Mode == MODE_IMMEDIATE {
		return a.V
	} else if a.Mode == MODE_POSITION {
		offset := toi(&a.V)
		return p.Memory[offset]
	} else if a.Mode == MODE_RELATIVE {
		offset := toi(&a.V)
		return p.Memory[p.RelativeBase + offset]
	}
	panic("invalid argument mode")
}

func (a Arg) PValue(p IntcodeCPU) *big.Int {
	if a.Mode == MODE_POSITION {
		offset := toi(&a.V)
		return &p.Memory[offset]
	} else if a.Mode == MODE_RELATIVE {
		offset := toi(&a.V)
		return &p.Memory[p.RelativeBase + offset]
	}
	panic("invalid pointer value")
}
