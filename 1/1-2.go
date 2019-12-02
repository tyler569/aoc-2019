
package main

import (
    "bufio"
    // "io/ioutil"
    "fmt"
    "log"
    // "io"
    "os"
    "strconv"
    "strings"
)

func read_ints_from_file(filename string) ([]int, error) {
    f, err := os.Open(filename)
    if (err != nil) {
        return nil, err
    }

    b := bufio.NewReader(f)
    int_array := []int{}

    for {
        line, err := b.ReadString('\n')
        if (err != nil) {
            break
        }
        line = strings.TrimSuffix(line, "\n")
        l_num, err := strconv.Atoi(line)
        if (err != nil) {
            break
        }

        int_array = append(int_array, l_num)
    }
    return int_array, nil
}

func weight_to_fuel(weight int) int {
    fuel := int(weight / 3) - 2
    if fuel < 0 {
        return 0
    }
    return fuel
}

func module_fuel(module int) int {
    if module == 0 {
        return 0
    }
    fuel := weight_to_fuel(module)
    fuel += module_fuel(fuel)
    return fuel
}

func main() {
    module_weights, err := read_ints_from_file("./input")
    if err != nil {
        log.Fatal(err)
    }

    total_fuel := 0
    for _, module := range module_weights {
        fuel := module_fuel(module)
        total_fuel += fuel
        fmt.Println("module:", module, "fuel:", fuel)
    }

    fmt.Println("Total fuel", total_fuel)
}

