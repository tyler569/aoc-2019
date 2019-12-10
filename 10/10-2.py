#!/usr/bin/env python3

from collections import defaultdict
import operator
import math

def delta(a, b):
    return (a[0] - b[0], a[1] - b[1])

def occluded(f, a):
    d = delta(a, f)
    step = math.gcd(d[0], d[1])

    xstep = d[0] // step
    ystep = d[1] // step

    x0 = a[0] + xstep
    y0 = a[1] + ystep
    while True:
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

candidate = (17, 14)
#candidate = (13, 11)

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
visible_cells = set(cells.keys()) - known_occluded

visible_asteroids = set()

asteroids = 0
for visible_cell in visible_cells:
    if visible_cell == candidate:
        continue
    if cells[visible_cell] == '#':
        visible_asteroids.add(visible_cell)

for i in range(len(rows)):
    for j in range(len(rows[0])):
        pos = (i, j)
        if pos == candidate:
            print("^", end="")
        elif cells[pos] == "#" and pos in known_occluded:
            print("!", end="")
        elif cells[pos] == "#":
            print("#", end="")
        elif pos in known_occluded:
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

assert (12, 11) in visible_asteroids

#print(visible_asteroids, len(visible_asteroids))

def angle_to(point, origin):
    d = delta(point, origin)
    angle = math.atan2(d[0], d[1])
    angle += math.pi/2
    if angle < 0:
        angle += 100
    return angle

angles = [(asteroid, angle_to(asteroid, candidate)) for asteroid in visible_asteroids]

angles.sort(key = lambda x: x[1])

print([a[0] for a in angles])
print(angles[199])

