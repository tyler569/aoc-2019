
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

func sliceMatchesRules(r []int) bool {
    last := -1
    double := false
    for _, v := range r {
        if v < last {
            return false
        }
        if v == last {
            double = true
        }
        last = v
    }
    return double
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

    count := 0
    for x := 136818; x <= 685979; x += 1 {
        if sliceMatchesRules(numberToSlice(x)) {
            count += 1
        }
    }
    fmt.Println("count:", count)
}

