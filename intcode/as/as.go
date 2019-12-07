package main

import (
	"fmt"

	"github.com/tyler569/aoc-2019/intcode/bits"
)

var opNames map[bits.Op]string = map[bits.Op]string{
	bits.PSOP_DATA:   "di",
	bits.OP_ADD:      "add",
	bits.OP_MULTIPLY: "mul",
	bits.OP_INPUT:    "in",
	bits.OP_OUTPUT:   "out",
	bits.OP_JUMP_T:   "jmp",
	bits.OP_JUMP_F:   "jmp",
	bits.OP_LESS:     "lt",
	bits.OP_EQUALS:   "eq",
	bits.OP_HALT:     "hlt",
}

func formatArg(a bits.Arg) string {
	if a.Mode == bits.MODE_IMMEDIATE {
		return fmt.Sprintf("%d", a.V)
	} else {
		return fmt.Sprintf("[%d]", a.V)
	}
}

func disas(o bits.Instruction) string {
	if o.Opcode == bits.PSOP_DATA {
		return fmt.Sprintf("di %#x", o.Args[0].V)
	}
	result := opNames[o.Opcode]
	for _, arg := range o.Args {
		result += " " + formatArg(arg)
	}
	return result
}

func parseAsmLine(line string) bits.Instruction {
	return bits.Instruction{}
}

func main() {
	/*
	program := bits.ReadProgram("input")
	instrs := bits.ScanInstructions(program)
	for _, i := range instrs {
		fmt.Println(disas(i))
	}
	*/

	file, err := os.Open("input.as")
	if err != nil {
		panic(err)
	}
}
