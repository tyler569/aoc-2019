#!/usr/bin/env python3

from collections import defaultdict
import operator
import math

def delta(a, b):
    return (a[0] - b[0], a[1] - b[1])

def occluded(f, a):
    # print("occluded({}, {})".format(f, a))
    d = delta(a, f)
    step = math.gcd(d[0], d[1])

    xstep = d[0] // step
    ystep = d[1] // step

    #print("step", xstep, ystep)

    x0 = a[0] + xstep
    y0 = a[1] + ystep
    while True:
        #print("  - {}".format((x0, y0)))
        yield (x0, y0)
        x0 += xstep
        y0 += ystep

text = ""
with open('input') as f:
    text = f.read()

rows = text.split("\n")
rows = list(filter(len, rows))

cells = {}
for i in range(len(rows)):
    for j in range(len(rows[0])):
        cells[(i, j)] = rows[i][j]

all_visible = {}

occluded_each = {}

for candidate in cells:
    if cells[candidate] != '#':
        continue
    known_occluded = set()
    for test in cells:
        if cells[test] != '#':
            continue
        if test == candidate:
            continue
        for x in occluded(candidate, test):
            if x not in cells:
                break
            known_occluded.add(x)
    occluded_each[candidate] = known_occluded
    visible_cells = set(cells.keys()) - known_occluded

    asteroids = 0
    for visible_cell in visible_cells:
        if visible_cell == candidate:
            continue
        if cells[visible_cell] == '#':
            asteroids += 1

    all_visible[candidate] = asteroids

print(all_visible)
best = max(all_visible.items(), key=operator.itemgetter(1))

for i in range(len(rows)):
    for j in range(len(rows[0])):
        pos = (i, j)
        if pos == best[0]:
            print("^", end="")
        elif cells[pos] == "#" and pos in occluded_each[best[0]]:
            print("!", end="")
        elif cells[pos] == "#":
            print("#", end="")
        elif pos in occluded_each[best[0]]:
            print("*", end="")
        elif cells[pos] == ".":
            print(" ", end="")
        else:
            print("?", end="")
        # elif pos in occluded_each[best[0]]:
        #     print("*", end="")
        # else:
        #     print(" ", end="")
    print()

print("best is {}, with {} visible".format(best[0], best[1]))

