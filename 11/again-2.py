#!/usr/bin/env python3

from collections import defaultdict
import subprocess
import sys

intcode = subprocess.Popen(["../intcode/intcode", "--program", "./input"],
        stdin=subprocess.PIPE, stdout=subprocess.PIPE)

input = intcode.stdin
output = intcode.stdout

paint, turn = 0, 0
next = 0

turn_state = 0
position = (0, 0)

def do_move(position, turn_state):
    x, y = position
    if turn_state == 0:
        return x, y+1
    elif turn_state == 1:
        return x+1, y
    elif turn_state == 2:
        return x, y-1
    elif turn_state == 3:
        return x-1, y

def do_turn(turn_state, turn_direction):
    if turn_direction == 0:
        turn_direction = -1
    new_turn = (turn_state + turn_direction) % 4
    if new_turn < 0:
        new_turn += 4
    return new_turn

hull = defaultdict(int)
hull[0, 0] = 1

input.write("{}\n".format(hull[position]).encode())
input.flush()

for line in output:
    if next == 0:
        paint = int(line)
        next = 1
        continue
    else:
        turn = int(line)
        next = 0

    print(paint, turn)

    hull[position] = paint

    turn_state = do_turn(turn_state, turn)
    position = do_move(position, turn_state)

    input.write("{}\n".format(hull[position]).encode())
    input.flush()

for y in range(10, -10, -1):
    for x in range(50):
        if hull[(x, y)]:
            print("#", end="")
        else:
            print(" ", end="")
    print()

