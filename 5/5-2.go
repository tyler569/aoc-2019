package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

type Mode int

const (
	MODE_POSITION  Mode = 0
	MODE_IMMEDIATE Mode = 1

	OP_ADD      = 1
	OP_MULTIPLY = 2
	OP_INPUT    = 3
	OP_OUTPUT   = 4
	OP_JUMP_T   = 5
	OP_JUMP_F   = 6
	OP_LESS     = 7
	OP_EQUALS   = 8
	OP_HALT     = 99
)

var modeStrings map[Mode]string = map[Mode]string{
	MODE_POSITION:  "pos",
	MODE_IMMEDIATE: "imm",
}

func (m Mode) String() string {
	return modeStrings[m]
}

var opStrings map[int]string = map[int]string{
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

type Instruction struct {
	opcode int
	arg1   Mode
	arg2   Mode
	arg3   Mode
}

func (m Instruction) String() string {
	return fmt.Sprintf("{%v %v %v %v}", opStrings[m.opcode],
		m.arg1, m.arg2, m.arg3)
}

func opModes(instr int) Instruction {
	opcode := instr % 100
	arg1_mode := Mode((instr / 100) % 10)
	arg2_mode := Mode((instr / 1000) % 10)
	arg3_mode := Mode((instr / 10000) % 10)

	return Instruction{opcode, arg1_mode, arg2_mode, arg3_mode}
}

func getArg(program []int, index int, mode Mode) int {
	if mode == MODE_POSITION {
		pos := program[index]
		return program[pos]
	} else if mode == MODE_IMMEDIATE {
		return program[index]
	} else {
		panic("Invalid parameter mode")
	}
}

func printInstruction(instruction []int) {
	m := opModes(instruction[0])
	/*
		if m.opcode == OP_ADD || m.opcode == OP_MULTIPLY {
			fmt.Println(m, instruction[1:4])
		} else if m.opcode == OP_INPUT || m.opcode == OP_OUTPUT {
			fmt.Println(m, instruction[1])
		} else if m.opcode == OP_HALT {
			fmt.Println("halt")
		} else if m.opcode == OP_JUMP_T || m.opcode == OP_JUMP_F {
			fmt.Println(m, instruction[1:3])
		} else if m.opcode == OP_LESS || m.opcode == OP_EQUALS {
			fmt.Println(m, instruction[1:4])
		} else {
			fmt.Println("invalid")
		}
	*/
	if m.opcode == OP_ADD || m.opcode == OP_MULTIPLY {
		fmt.Println(instruction[0:4])
	} else if m.opcode == OP_INPUT || m.opcode == OP_OUTPUT {
		fmt.Println(instruction[0:1])
	} else if m.opcode == OP_HALT {
		fmt.Println("halt")
	} else if m.opcode == OP_JUMP_T || m.opcode == OP_JUMP_F {
		fmt.Println(instruction[0:3])
	} else if m.opcode == OP_LESS || m.opcode == OP_EQUALS {
		fmt.Println(instruction[0:4])
	} else {
		fmt.Println("invalid -", instruction[0:1])
	}
}

func executeProgram(program Program) Program {
	index := 0

	for {
		// printInstruction(program[index:])

		instruction := opModes(program[index])
		opcode := instruction.opcode

		if opcode == OP_HALT {
			break
		} else if opcode == OP_ADD || opcode == OP_MULTIPLY {
			out_index := program[index+3]

			in1 := getArg(program, index+1, instruction.arg1)
			in2 := getArg(program, index+2, instruction.arg2)

			var out int
			if opcode == OP_ADD {
				out = in1 + in2
			} else if opcode == OP_MULTIPLY {
				out = in1 * in2
			}
			program[out_index] = out

			index += 4
		} else if opcode == OP_INPUT {
			input_index := program[index+1]
			fmt.Printf("input: ")

			var i int
			_, err := fmt.Scanf("%d", &i)
			if err != nil {
				panic(err)
			}
			program[input_index] = i

			index += 2
		} else if opcode == OP_OUTPUT {
			output_val := getArg(program, index+1, instruction.arg1)
			fmt.Printf("output: %d\n", output_val)

			index += 2
		} else if opcode == OP_JUMP_F || opcode == OP_JUMP_T {
			opt := getArg(program, index+1, instruction.arg1)
			target := getArg(program, index+2, instruction.arg2)

			if (opcode == OP_JUMP_F && opt == 0) || (opcode == OP_JUMP_T && opt != 0) {
				// jump to target
				index = target
			} else {
				index += 3
			}
		} else if opcode == OP_LESS || opcode == OP_EQUALS {
			p1 := getArg(program, index+1, instruction.arg1)
			p2 := getArg(program, index+2, instruction.arg2)
			out_index := program[index+3]

			if (opcode == OP_LESS && p1 < p2) || (opcode == OP_EQUALS && p1 == p2) {
				program[out_index] = 1
			} else {
				program[out_index] = 0
			}

			index += 4
		} else {
			// error print routine

			bot := max(0, index-5)
			top := min(index+6, len(program))

			for j := bot; j < top; j += 1 {
				var ic string
				if j == index {
					ic = "*"
				} else {
					ic = ""
				}

				fmt.Printf("%s%d%[1]s ", ic, program[j])
			}
			fmt.Printf("\n")

			fmt.Println("index:", index)
			panic("Invalid op at index")
		}
		// fmt.Println(program)
	}

	return program
}

func main() {
	program := readProgram("input")

	// fmt.Println(program)

	executeProgram(program)
}
