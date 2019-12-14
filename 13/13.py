#!/usr/bin/env python3

from collections import defaultdict
import subprocess
import sys

intcode = subprocess.Popen(["../intcode/intcode", "--program", "./input"],
        stdin=subprocess.PIPE, stdout=subprocess.PIPE)

input = intcode.stdin
output = intcode.stdout

lines = output.read().decode().split("\n")
lines = filter(lambda a: a, lines)
lines = list(map(int, lines))

chunks = [lines[a:a+3] for a in range(0, len(lines), 3)]
chunks = list(map(lambda c: ((c[0], c[1]), c[2]), chunks))

poss = {}

for c in chunks:
    if c[0] in poss:
        print("overwriting", poss[c[0]])
    poss[c[0]] = c[1]

i = 0

for c, v in chunks:
    if v == 2:
        i += 1

print(i)

