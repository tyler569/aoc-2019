#!/usr/bin/env python3

import sys
import math
from collections import defaultdict

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

def make_graph(deps):
    print('digraph deps {')
    for k in deps:
        q_k, reqs = deps[k]
        for q_r, r in reqs:
            print(r, '->', k, '[ label="{} -> {}" ];'.format(q_r, q_k))
    print('}')

def resolve(times, request, have):
    print('xxx', times, request, have)
    ore_need = 0
    for r in request:
        if r[1] == 'ORE':
            ore_need += r[0] * times
            continue
        need = r[0] * times
        qty = 0
        if r[1] in have and have[r[1]] > 0:
            take = min(have[r[1]], need)
            print('take', take, r[1])
            have[r[1]] -= take
            need -= take
        
        if need > 0:
            makes, dep = deps[r[1]]
            cycles = math.ceil(need / makes)
            qty = cycles * makes
            print('make', qty, r[1])
            ore_need += resolve(cycles, dep, have)

        if qty > need:
            if r[1] in have:
                have[r[1]] += qty - need
            else:
                have[r[1]] = qty - need
    return ore_need

print(resolve(1, [(1, 'FUEL')], {}))#defaultdict(int)))

# make_graph(deps)
