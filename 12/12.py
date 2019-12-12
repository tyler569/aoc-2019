#!/usr/bin/env python3

import copy
import itertools
import sys

moons = []

for line in sys.stdin:
    if line.strip() == "":
        continue
    moonpos = line.strip().replace("<", "").replace(">", "")
    moonpos = moonpos.split(",")
    moonpos = list(map(lambda s: int(s.split("=")[1]), moonpos))
    moons.append((moonpos, [0, 0, 0]))

for moon in moons:
    print(moon)

def gravity(moon1, moon2):
    delta1 = [0, 0, 0]
    delta2 = [0, 0, 0]

    for axis in range(3):
        if moon1[0][axis] > moon2[0][axis]:
            moon1[1][axis] -= 1
            moon2[1][axis] += 1
        elif moon1[0][axis] < moon2[0][axis]:
            moon1[1][axis] += 1
            moon2[1][axis] -= 1

def energy(moon):
    return sum(map(abs, moon[0])) * sum(map(abs, moon[1]))

def total_energy(moons):
    return sum(map(energy, moons))

previousX = set()
previousY = set()
previousZ = set()

for i in range(10000000):
    # print("step", i)
    for pair in itertools.combinations(moons, 2):
        moon1, moon2 = pair
        gravity(moon1, moon2)

    for moon in moons:
        for axis in range(3):
            moon[0][axis] += moon[1][axis]
        # print(moon)

    pX = tuple(map(lambda a: (a[0][0], a[1][0]), moons))
    pY = tuple(map(lambda a: (a[0][1], a[1][1]), moons))
    pZ = tuple(map(lambda a: (a[0][2], a[1][2]), moons))

    if previousX and pX in previousX:
        print("x:", i)
        previousX = None

    if previousY and pY in previousY:
        print("x:", i)
        previousY = None

    if previousZ and pZ in previousZ:
        print("x:", i)
        previousZ = None

    previousX is not None and previousX.add(pX)
    previousY is not None and previousY.add(pY)
    previousZ is not None and previousZ.add(pZ)

    if not previousX and not previousY and not previousZ:
        break


