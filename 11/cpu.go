package main

import (
	"fmt"
	"os"
	// "runtime"
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
			fmt.Println("INPUT INSTRUCTION")
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

type Point struct {
	x, y int
}

type Color int

const (
	BLACK Color = 0
	WHITE Color = 1
)


type Direction int

const (
	UP Direction = 0
	RIGHT Direction = 1
	DOWN Direction = 2
	LEFT Direction = 3
)

const (
	TURN_LEFT = 0
	TURN_RIGHT = 1
)

type Painter struct {
	dir Direction
	at Point
}

func (p *Painter) Turn(turn_type int) {
	if turn_type == TURN_LEFT {
		p.dir -= Direction(1)
	} else if turn_type == TURN_RIGHT {
		p.dir += Direction(1)
	} else {
		panic(turn_type)
	}

	p.dir %= Direction(4)
	if p.dir < Direction(0) {
		p.dir += Direction(4)
	}
}

func (p *Painter) Move() {
	if p.dir == UP {
		p.at.y += 1
	} else if p.dir == RIGHT {
		p.at.x += 1
	} else if p.dir == DOWN {
		p.at.y -= 1
	} else if p.dir == LEFT {
		p.at.x -= 1
	} else {
		fmt.Println(p.dir)
		panic(p.dir)
	}
}

func (d Direction) String() string {
	switch (d) {
	case UP:
		return "up"
	case DOWN:
		return "down"
	case LEFT:
		return "left"
	case RIGHT:
		return "right"
	default:
		panic(d)
	}
}

func main() {
	file, err := os.Open("input")
	if err != nil {
		panic(err)
	}
	program := ReadProgram(file)

	// fmt.Println(program)

	/*
		instrs := ScanInstructions(program)
		for _, v := range instrs {
			fmt.Println(v)
		}
	*/
	input := make(chan int)
	output := make(chan int)

	memory := make([]int, 10000)
	copy(memory, program)

	cpu := IntcodeCPU{
		Memory:       memory,
		RelativeBase: 0,
		ProgramCounter: 0,
	}

	go ExecuteProgram(cpu, input, output)

	var wg sync.WaitGroup
	wg.Add(1)

	var painter Painter
	painter.dir = UP

	Hull := make(map[Point]Color)
	Hull[Point{0, 0}] = WHITE
	HullMtx := sync.Mutex{}

	go func() {
		var paint, turn int
		var next int
		for c := range output {
			if next == 0 {
				HullMtx.Lock()
				paint = c
				next = 1
			} else if next == 1 {
				turn = c

				fmt.Println("pair:", paint, turn)
				if turn == 0 {
					fmt.Println("turing left")
				} else {
					fmt.Println("turning right")
				}

				Hull[painter.at] = Color(paint)
				// fmt.Printf("Painted color %v at %v\n", Color(paint), painter.at)
				// fmt.Printf("Turned to face %v\n", painter.dir)

				painter.Turn(turn)
				painter.Move()

				next = 0
				HullMtx.Unlock()
			}
		}
		wg.Done()
	}()

	go func() {
		for {
			input <- 0
			HullMtx.Lock()
			color, _ := Hull[painter.at] // no ok because 0 is correct
			fmt.Printf("Read color %v from %v\n", color, painter.at)
			input <- int(color)
			HullMtx.Unlock()
		}
	}()

	wg.Wait()

	fmt.Println(Hull)
	fmt.Println(len(Hull))

	for y := 10; y > -10; y-- {
		for x := 00; x < 50; x++ {
			color, ok := Hull[Point{x, y}]
			if !ok || color == BLACK {
				fmt.Print(".")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
}

