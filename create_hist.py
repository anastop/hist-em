import numpy as np
import subprocess
import os
import subprocess
import signal
import sys

def print_hist():

    with open('output.txt') as f:
        lines = f.readlines()

    start_hist = False
    values = []
    for l in lines:
        if 'Total' in l:
            break
        if start_hist:
            values.append(l)
        if 'Histogram' in l:
            start_hist = True

    final_array = []
    for val in values:
        _l = val.split(' ')
        tm, occs = int(_l[0]), int(_l[1])
        if tm == 0 or occs == 0:
           continue
        final_array.extend([tm for _ in range(occs)])
    count, division = np.histogram(final_array, bins=20, range=(0.0, 100.0))
    for idx, _ in enumerate(count):
        print int(division[idx]), count[idx]

	
print_hist()



