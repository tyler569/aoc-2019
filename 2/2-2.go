
package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    "strconv"
)

func read_ints_from_file(filename string) ([]int, error) {
    f, err := os.Open(filename)
    if (err != nil) {
        return nil, err
    }

    b := bufio.NewReader(f)
    int_array := []int{}

    for {
        line, err := b.ReadString(',')
        if (err != nil) {
            break
        }
        line = strings.TrimSuffix(line, ",")
        l_num, err := strconv.Atoi(line)
        if (err != nil) {
            break
        }

        int_array = append(int_array, l_num)
    }
    return int_array, nil
}

func execute_program(program []int) []int {
    index := 0

    for {
        if program[index] == 99 { // halt
            break
        } else if program[index] == 1 { // +
            in1_index := program[index + 1]
            in2_index := program[index + 2]
            out_index := program[index + 3]

            in1 := program[in1_index]
            in2 := program[in2_index]

            out := in1 + in2
            program[out_index] = out

            index += 4
        } else if program[index] == 2 { // *
            in1_index := program[index + 1]
            in2_index := program[index + 2]
            out_index := program[index + 3]

            in1 := program[in1_index]
            in2 := program[in2_index]

            out := in1 * in2
            program[out_index] = out

            index += 4
        }
    }

    return program
}

func main() {
    original_program, err := read_ints_from_file("input")
    if err != nil {
        log.Fatal(err)
    }

    for i := 0; i < 100; i++ {
        for j := 0; j < 100; j++ {
            program := make([]int, len(original_program))
            copy(program, original_program)
            // setup
            program[1] = i
            program[2] = j

            execute_program(program)

            if program[0] == 19690720 {
                fmt.Println(i, j)
                break
            }
        }
    }
}

