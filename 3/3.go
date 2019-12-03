
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

type Turn struct {
	direction Direction
	distance int
}

type Point struct {
    x, y int
}

func manhattanDistanceTo0(p Point) int {
    return abs(p.x) + abs(p.y)
}

type Segment struct {
    s, e Point
}

type Wire struct {
	turns []Turn
    segments []Segment
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

func (t Turn) String() string {
    //return fmt.Sprintf("{%v %d}", t.direction, t.distance)
    return "Turn{...}"
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

        w.segments = append(w.segments, Segment{p_prev, p_next})
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

func segmentIntersection(s1, s2 Segment) (Point, error) {
    d1 := absSegmentDirection(s1)
    d2 := absSegmentDirection(s2)
    if d1 == d2 {
        // parllel, no interseciton
        return Point{}, errors.New("No intersection")
    }

    if d1 == RIGHT {
        return segmentIntersection(s2, s1)
    }

    if between(s2.s.y, s1.s.y, s1.e.y) && between(s1.s.x, s2.s.x, s2.e.x) {
        // we have a winner
        return Point{s1.s.x, s2.s.y}, nil
    }

    return Point{}, errors.New("No intersection")
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

    var intersections []Point

    for _, s1 := range w1.segments {
        for _, s2 := range w2.segments {
            intersection, err := segmentIntersection(s1, s2)
            if err != nil {
                continue
            }
            if intersection.x == 0 && intersection.y == 0 {
                continue
            }
            intersections = append(intersections, intersection)
        }
    }

    fmt.Println(intersections)

    minimum := 100000
    for _, intersection := range intersections {
        distance := manhattanDistanceTo0(intersection)
        if distance < minimum {
            minimum = distance
        }
    }

    fmt.Println("minmum distance is", minimum)
}

