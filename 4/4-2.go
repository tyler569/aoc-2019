
package main

import (
    "fmt"
)

func numberToSlice(x int) []int {
    result := []int{
        x / 100000,
        x % 100000 / 10000,
        x % 10000 / 1000,
        x % 1000 / 100,
        x % 100 / 10,
        x % 10,
    }
    return result
}

func rotate(r []int) []int {
    return append(r[1:], r[0])
}

func sliceMatchesRules(r []int) bool {
    last := -1

    // scan for increase
    for _, v := range r {
        if v < last {
            return false
        }
        last = v
    }

    // a double == 2 long
    r_r := r
    for i := 0; i < len(r); i += 1 {
        _ = i
        if r_r[0] != r_r[1] &&
           r_r[1] == r_r[2] &&
           r_r[2] != r_r[3] {

            return true
        }
        r_r = rotate(r_r)
    }

    return false
}

func printTest(num int) {
    x := numberToSlice(num)
    y := sliceMatchesRules(x)
    fmt.Println(x, y)
}

func main() {
    printTest(111111)
    printTest(223450)
    printTest(123789)
    printTest(111123)
    printTest(112233)
    printTest(123444)
    printTest(111122)

    count := 0
    for x := 136818; x <= 685979; x += 1 {
        if sliceMatchesRules(numberToSlice(x)) {
            count += 1
        }
    }
    fmt.Println("count:", count)
}

