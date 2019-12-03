
package main

import (
    "bufio"
    "errors"
    "fmt"
    "io"
	"log"
	"os"
    "strings"
    "strconv"
)

/*
NB: this file makes the strong assumption that Segments are always exactly
either left/right or up/down. Don't try to extend this for diagonal without
auditing everything
*/

func abs(a int) int {
    if a < 0 {
        return -a
    }
    return a
}

type Direction int

const (
	RIGHT Direction = iota
	LEFT
	UP
	DOWN
)

var dir_names = map[Direction]string{
    RIGHT: "right",
    LEFT: "left",
    UP: "up",
    DOWN: "down",
}

type Turn struct {
	direction Direction
	distance int
}

func (t Turn) String() string {
    return fmt.Sprintf("{%s %d}", dir_names[t.direction], t.distance)
}

type Point struct {
    x, y int
}

func manhattanDistanceTo0(p Point) int {
    return abs(p.x) + abs(p.y)
}

type Segment struct {
    s, e Point
    length int
    steps  int
}

func (s Segment) String() string {
    return fmt.Sprintf("{%v -> %v, len: %d, steps: %d}",
        s.s, s.e, s.length, s.steps)
}

type Wire struct {
	turns    []Turn
    segments []Segment
}

func (w Wire) String() string {
    b := strings.Builder{}
    b.WriteString("{ turns: {")
    for _, t := range w.turns {
        b.WriteString(fmt.Sprint(t, ","))
    }
    b.WriteString("}, segments: {")
    for _, s := range w.segments {
        b.WriteString(fmt.Sprint(s, ",\n"))
    }
    b.WriteString("}")
    return b.String()
}

type Intersection struct {
    p      Point
    x1, x2 int // steps for each wire to get here
}

func parseTurn(input string) (Turn, error) {
    var t Turn
    switch input[0] {
    case 'R':
        t.direction = RIGHT
    case 'L':
        t.direction = LEFT
    case 'U':
        t.direction = UP
    case 'D':
        t.direction = DOWN
    default:
        return t, errors.New("Invalid Direction")
    }

    distance, err := strconv.Atoi(input[1:])
    if err != nil {
        return t, err
    }

    t.distance = distance
    return t, nil
}

func parseWire(input string) (Wire, error) {
    string_read := strings.NewReader(input)
    s := bufio.NewReader(string_read)

    w := Wire{}

    for {
        var t Turn
        eof := false
        segment, err := s.ReadString(',')
        if err == nil {
            segment = strings.TrimSuffix(segment, ",")
            t, err = parseTurn(segment)
            if err != nil {
                return w, err
            }
        } else if err == io.EOF {
            eof = true
            segment = strings.TrimSuffix(segment, "\n")
            t, err = parseTurn(segment)
            if err != nil {
                return w, err
            }
        } else {
            return w, err
        }

        w.turns = append(w.turns, t)

        if eof {
            break
        }
    }

    p0 := Point{0, 0}
    p_prev := p0

    total_steps := 0


    for _, t := range w.turns {
        p_next := p_prev
        switch t.direction {
        case RIGHT:
            p_next.x += t.distance
        case LEFT:
            p_next.x -= t.distance
        case UP:
            p_next.y += t.distance
        case DOWN:
            p_next.y -= t.distance
        }

        w.segments = append(w.segments,
            Segment{p_prev, p_next, t.distance, total_steps})

        total_steps += t.distance
        p_prev = p_next
    }

    return w, nil
}

func absSegmentDirection(s Segment) Direction {
    if s.s.x == s.e.x {
        return UP
    }
    return RIGHT
}

func between(i, a, b int) bool {
    between := ((i >= a) && (i <= b)) || ((i >= b) && (i <= a))

    return between
}

func segmentIntersection(s1, s2 Segment) (Intersection, error) {
    d1 := absSegmentDirection(s1)
    d2 := absSegmentDirection(s2)
    if d1 == d2 {
        // parllel, no interseciton
        return Intersection{}, errors.New("No intersection")
    }

    if d1 == RIGHT {
        return segmentIntersection(s2, s1)
    }

    if between(s2.s.y, s1.s.y, s1.e.y) && between(s1.s.x, s2.s.x, s2.e.x) {
        // we have a winner
        p := Point{s1.s.x, s2.s.y}

        s1_delta := abs(s2.s.y - s1.s.y)
        s2_delta := abs(s2.s.x - s1.s.x)

        w1_steps := s1.steps + s1_delta
        w2_steps := s2.steps + s2_delta

        return Intersection{p, w1_steps, w2_steps}, nil
    }

    return Intersection{}, errors.New("No insersection")
}

func main() {
	input, err := os.Open("input")
	if err != nil {
		log.Fatal(err)
	}

    b := bufio.NewReader(input)
    s, err := b.ReadString('\n')
    if err != nil {
        log.Fatal(err)
    }
    w1, err := parseWire(s)
    if err != nil {
        log.Fatal(err)
    }

    s, err = b.ReadString('\n')
    if err != nil {
        log.Fatal(err)
    }
    w2, err := parseWire(s)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(w1)
    fmt.Println(w2)
    var intersections []Intersection

    for _, s1 := range w1.segments {
        for _, s2 := range w2.segments {
            intersection, err := segmentIntersection(s1, s2)
            if err != nil {
                continue
            }
            if intersection.p.x == 0 && intersection.p.y == 0 {
                continue
            }
            intersections = append(intersections, intersection)
        }
    }

    fmt.Println(intersections)

    minimum := 100000
    for _, intersection := range intersections {
        total_steps := intersection.x1 + intersection.x2
        if total_steps < minimum {
            minimum = total_steps
        }
    }

    fmt.Println("minimum combined steps is", minimum)
}

