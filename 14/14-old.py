#!/usr/bin/env python3

import sys
import math

deps = {}
for line in sys.stdin:
    if line.strip() == "":
        continue
    input, output = line.split("=>")
    inputs = input.split(",")
    
    input_q = []
    for i in inputs:
        num, thing = i.split()
        input_q.append((int(num), thing))

    output = output.split()
    output_num = int(output[0])
    output_thing = output[1]

    deps[output_thing] = (output_num, input_q)

def ore_for(v0):
    print('ore_for(' + str(v0) + ')')
    q, x = v0
    if x == 'ORE':
        return q

    makes, ds = deps[x]

    if q % makes != 0:
        print('inexact -', q, 'does not divide', makes)

    times = math.ceil(q / makes)

    print('making', times, 'times', x)

    return times * sum(map(ore_for, ds))

print(ore_for((1, 'FUEL')))

